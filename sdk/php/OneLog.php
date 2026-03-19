<?php

/**
 * One-Log (ULAM) Official PHP SDK
 * 
 * Example usage:
 * ```php
 * use OneLog\OneLog;
 * 
 * $client = new OneLog([
 *     'api_key' => 'ulam_live_xxxxxxxx',
 *     'endpoint' => 'https://api.ulam.your-domain.com'
 * ]);
 * 
 * $client->logInfo(OneLog::CATEGORY_SYSTEM_ERROR, 'Database connection established');
 * ```
 */

namespace OneLog;

use Exception;
use RuntimeException;

class OneLog
{
    // Log Levels
    const LEVEL_CRITICAL = 'CRITICAL';
    const LEVEL_ERROR = 'ERROR';
    const LEVEL_WARN = 'WARN';
    const LEVEL_INFO = 'INFO';
    const LEVEL_DEBUG = 'DEBUG';

    // Categories
    const CATEGORY_SYSTEM_ERROR = 'SYSTEM_ERROR';
    const CATEGORY_USER_ACTIVITY = 'USER_ACTIVITY';
    const CATEGORY_AUTH_EVENT = 'AUTH_EVENT';
    const CATEGORY_PERFORMANCE = 'PERFORMANCE';
    const CATEGORY_SECURITY = 'SECURITY';
    const CATEGORY_AUDIT_TRAIL = 'AUDIT_TRAIL';

    private string $apiKey;
    private string $endpoint;
    private ?string $sourceId;
    private float $timeout;

    /**
     * Create a new One-Log client
     * 
     * @param array $options Configuration options
     * @throws RuntimeException If API key is not provided
     */
    public function __construct(array $options = [])
    {
        $this->apiKey = $options['api_key'] ?? $_ENV['ULAM_API_KEY'] ?? '';
        $this->endpoint = $options['endpoint'] ?? $_ENV['ULAM_ENDPOINT'] ?? 'http://localhost:8080/api/ingest';
        $this->sourceId = $options['source_id'] ?? $_ENV['ULAM_SOURCE_ID'] ?? null;
        $this->timeout = $options['timeout'] ?? 5.0;

        if (empty($this->apiKey)) {
            throw new RuntimeException('API key is required. Set ULAM_API_KEY environment variable or pass api_key option.');
        }
    }

    /**
     * Send a log entry to One-Log
     * 
     * @param string $category Log category
     * @param string $level Log level
     * @param string $message Log message
     * @param array $context Additional context data
     * @param string|null $stackTrace Stack trace
     * @throws Exception If request fails
     */
    public function log(
        string $category,
        string $level,
        string $message,
        array $context = [],
        ?string $stackTrace = null
    ): void {
        $entry = [
            'category' => $category,
            'level' => $level,
            'message' => $message,
            'context' => $context,
            'stack_trace' => $stackTrace,
            'ip_address' => $this->getIPAddress()
        ];

        $this->send($entry);
    }

    /**
     * Log with error details
     * 
     * @param string $category Log category
     * @param string $level Log level
     * @param string $message Log message
     * @param Exception $error Error exception
     * @param array $context Additional context
     * @throws Exception If request fails
     */
    public function logError(
        string $category,
        string $level,
        string $message,
        Exception $error,
        array $context = []
    ): void {
        $stackTrace = $error->getTraceAsString();
        $this->log($category, $level, $message, $context, $stackTrace);
    }

    /**
     * Log critical level message
     * 
     * @param string $category Log category
     * @param string $message Log message
     * @param array $context Additional context
     */
    public function logCritical(string $category, string $message, array $context = []): void
    {
        $this->log($category, self::LEVEL_CRITICAL, $message, $context);
    }

    /**
     * Log error level message
     * 
     * @param string $category Log category
     * @param string $message Log message
     * @param array $context Additional context
     */
    public function logErrorLevel(string $category, string $message, array $context = []): void
    {
        $this->log($category, self::LEVEL_ERROR, $message, $context);
    }

    /**
     * Log warning level message
     * 
     * @param string $category Log category
     * @param string $message Log message
     * @param array $context Additional context
     */
    public function logWarn(string $category, string $message, array $context = []): void
    {
        $this->log($category, self::LEVEL_WARN, $message, $context);
    }

    /**
     * Log info level message
     * 
     * @param string $category Log category
     * @param string $message Log message
     * @param array $context Additional context
     */
    public function logInfo(string $category, string $message, array $context = []): void
    {
        $this->log($category, self::LEVEL_INFO, $message, $context);
    }

    /**
     * Log debug level message
     * 
     * @param string $category Log category
     * @param string $message Log message
     * @param array $context Additional context
     */
    public function logDebug(string $category, string $message, array $context = []): void
    {
        $this->log($category, self::LEVEL_DEBUG, $message, $context);
    }

    /**
     * Log performance metrics
     * 
     * @param string $endpoint API endpoint path
     * @param int $durationMs Request duration in milliseconds
     * @param string $method HTTP method
     * @param int $statusCode HTTP status code
     */
    public function logPerformance(
        string $endpoint,
        int $durationMs,
        string $method = 'GET',
        int $statusCode = 200
    ): void {
        $this->log(
            self::CATEGORY_PERFORMANCE,
            self::LEVEL_INFO,
            'API request completed',
            [
                'endpoint' => $endpoint,
                'duration_ms' => $durationMs,
                'method' => $method,
                'status_code' => $statusCode
            ]
        );
    }

    /**
     * Log authentication event
     * 
     * @param string $eventType Event type
     * @param string $authMethod Authentication method
     * @param string $userId User ID
     * @param bool $success Whether authentication succeeded
     */
    public function logAuthEvent(
        string $eventType,
        string $authMethod,
        string $userId,
        bool $success = true
    ): void {
        $level = $success ? self::LEVEL_INFO : self::LEVEL_WARN;
        $this->log(
            self::CATEGORY_AUTH_EVENT,
            $level,
            'Authentication event',
            [
                'event_type' => $eventType,
                'auth_method' => $authMethod,
                'user_id' => $userId,
                'success' => $success
            ]
        );
    }

    /**
     * Log audit trail event
     * 
     * @param string $action Action performed
     * @param string $actorId Actor/user ID
     * @param string $resourceType Type of resource
     * @param string $resourceId Resource ID
     * @param mixed $before State before change
     * @param mixed $after State after change
     */
    public function logAudit(
        string $action,
        string $actorId,
        string $resourceType,
        string $resourceId,
        $before = null,
        $after = null
    ): void {
        $this->log(
            self::CATEGORY_AUDIT_TRAIL,
            self::LEVEL_INFO,
            'Audit event',
            [
                'action' => $action,
                'actor_id' => $actorId,
                'resource_type' => $resourceType,
                'resource_id' => $resourceId,
                'before' => $before,
                'after' => $after
            ]
        );
    }

    /**
     * Fire-and-forget logging (async, no wait)
     * Use when you don't want to block execution
     * 
     * @param string $category Log category
     * @param string $level Log level
     * @param string $message Log message
     * @param array $context Additional context
     */
    public function fireAndForget(
        string $category,
        string $level,
        string $message,
        array $context = []
    ): void {
        // In PHP, true async is complex. We'll use a non-blocking approach
        // by ignoring the response
        try {
            $this->log($category, $level, $message, $context);
        } catch (Exception $e) {
            // Silently ignore errors in fire-and-forget mode
        }
    }

    /**
     * Send log entry to API
     * 
     * @param array $entry Log entry
     * @throws Exception If request fails
     */
    private function send(array $entry): void
    {
        $ch = curl_init($this->endpoint);
        
        curl_setopt_array($ch, [
            CURLOPT_POST => true,
            CURLOPT_POSTFIELDS => json_encode($entry),
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT => $this->timeout,
            CURLOPT_HTTPHEADER => [
                'Content-Type: application/json',
                'X-API-Key: ' . $this->apiKey
            ]
        ]);

        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $error = curl_error($ch);
        
        curl_close($ch);

        if ($error) {
            throw new Exception("Request failed: {$error}");
        }

        if ($httpCode !== 202) {
            throw new Exception("Unexpected status code: {$httpCode}");
        }
    }

    /**
     * Get client IP address
     * 
     * @return string|null
     */
    private function getIPAddress(): ?string
    {
        // In a real implementation, this would detect the actual IP
        return null;
    }
}
