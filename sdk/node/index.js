/**
 * One-Log (ULAM) Official Node.js SDK
 * 
 * @example
 * const OneLog = require('@onelog/sdk');
 * 
 * const client = new OneLog({
 *   apiKey: 'ulam_live_xxxxxxxx',
 *   endpoint: 'https://api.ulam.your-domain.com'
 * });
 * 
 * await client.logInfo('SYSTEM_ERROR', 'Database connection established');
 */

class OneLog {
  /**
   * Create a new One-Log client
   * @param {Object} options - Configuration options
   * @param {string} options.apiKey - Your ULAM API key
   * @param {string} options.endpoint - API endpoint URL
   * @param {string} options.sourceId - Source identifier
   * @param {number} options.timeout - Request timeout in ms (default: 5000)
   */
  constructor(options = {}) {
    this.apiKey = options.apiKey || process.env.ULAM_API_KEY;
    this.endpoint = options.endpoint || process.env.ULAM_ENDPOINT || 'http://localhost:8080/api/ingest';
    this.sourceId = options.sourceId || process.env.ULAM_SOURCE_ID;
    this.timeout = options.timeout || 5000;

    if (!this.apiKey) {
      throw new Error('API key is required. Set ULAM_API_KEY environment variable or pass apiKey option.');
    }
  }

  /**
   * Log levels
   */
  static Levels = {
    CRITICAL: 'CRITICAL',
    ERROR: 'ERROR',
    WARN: 'WARN',
    INFO: 'INFO',
    DEBUG: 'DEBUG'
  };

  /**
   * Categories
   */
  static Categories = {
    SYSTEM_ERROR: 'SYSTEM_ERROR',
    USER_ACTIVITY: 'USER_ACTIVITY',
    AUTH_EVENT: 'AUTH_EVENT',
    PERFORMANCE: 'PERFORMANCE',
    SECURITY: 'SECURITY',
    AUDIT_TRAIL: 'AUDIT_TRAIL'
  };

  /**
   * Send a log entry to One-Log
   * @param {string} category - Log category
   * @param {string} level - Log level
   * @param {string} message - Log message
   * @param {Object} context - Additional context data
   * @param {string} stackTrace - Stack trace (optional)
   * @returns {Promise<void>}
   */
  async log(category, level, message, context = {}, stackTrace = null) {
    const entry = {
      category,
      level,
      message,
      context,
      stack_trace: stackTrace,
      ip_address: this._getIPAddress()
    };

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), this.timeout);

      const response = await fetch(this.endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': this.apiKey
        },
        body: JSON.stringify(entry),
        signal: controller.signal
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    } catch (error) {
      throw new Error(`Failed to send log: ${error.message}`);
    }
  }

  /**
   * Log with error details
   * @param {string} category - Log category
   * @param {string} level - Log level
   * @param {string} message - Log message
   * @param {Error} error - Error object
   * @param {Object} context - Additional context
   * @returns {Promise<void>}
   */
  async logError(category, level, message, error, context = {}) {
    const stackTrace = error ? error.stack : null;
    return this.log(category, level, message, context, stackTrace);
  }

  /**
   * Log critical level message
   * @param {string} category - Log category
   * @param {string} message - Log message
   * @param {Object} context - Additional context
   * @returns {Promise<void>}
   */
  async logCritical(category, message, context = {}) {
    return this.log(category, OneLog.Levels.CRITICAL, message, context);
  }

  /**
   * Log error level message
   * @param {string} category - Log category
   * @param {string} message - Log message
   * @param {Object} context - Additional context
   * @returns {Promise<void>}
   */
  async logErrorLevel(category, message, context = {}) {
    return this.log(category, OneLog.Levels.ERROR, message, context);
  }

  /**
   * Log warning level message
   * @param {string} category - Log category
   * @param {string} message - Log message
   * @param {Object} context - Additional context
   * @returns {Promise<void>}
   */
  async logWarn(category, message, context = {}) {
    return this.log(category, OneLog.Levels.WARN, message, context);
  }

  /**
   * Log info level message
   * @param {string} category - Log category
   * @param {string} message - Log message
   * @param {Object} context - Additional context
   * @returns {Promise<void>}
   */
  async logInfo(category, message, context = {}) {
    return this.log(category, OneLog.Levels.INFO, message, context);
  }

  /**
   * Log debug level message
   * @param {string} category - Log category
   * @param {string} message - Log message
   * @param {Object} context - Additional context
   * @returns {Promise<void>}
   */
  async logDebug(category, message, context = {}) {
    return this.log(category, OneLog.Levels.DEBUG, message, context);
  }

  /**
   * Log performance metrics
   * @param {string} endpoint - API endpoint path
   * @param {number} durationMs - Request duration in milliseconds
   * @param {string} method - HTTP method
   * @param {number} statusCode - HTTP status code
   * @returns {Promise<void>}
   */
  async logPerformance(endpoint, durationMs, method = 'GET', statusCode = 200) {
    return this.log(
      OneLog.Categories.PERFORMANCE,
      OneLog.Levels.INFO,
      'API request completed',
      { endpoint, duration_ms: durationMs, method, status_code: statusCode }
    );
  }

  /**
   * Log authentication event
   * @param {string} eventType - Event type (e.g., 'login_success', 'login_failed')
   * @param {string} authMethod - Authentication method
   * @param {string} userId - User ID
   * @param {boolean} success - Whether authentication succeeded
   * @returns {Promise<void>}
   */
  async logAuthEvent(eventType, authMethod, userId, success = true) {
    const level = success ? OneLog.Levels.INFO : OneLog.Levels.WARN;
    return this.log(
      OneLog.Categories.AUTH_EVENT,
      level,
      'Authentication event',
      { event_type: eventType, auth_method: authMethod, user_id: userId, success }
    );
  }

  /**
   * Log audit trail event
   * @param {string} action - Action performed
   * @param {string} actorId - Actor/user ID
   * @param {string} resourceType - Type of resource
   * @param {string} resourceId - Resource ID
   * @param {*} before - State before change
   * @param {*} after - State after change
   * @returns {Promise<void>}
   */
  async logAudit(action, actorId, resourceType, resourceId, before = null, after = null) {
    return this.log(
      OneLog.Categories.AUDIT_TRAIL,
      OneLog.Levels.INFO,
      'Audit event',
      { action, actor_id: actorId, resource_type: resourceType, resource_id: resourceId, before, after }
    );
  }

  /**
   * Fire-and-forget logging (async, no wait)
   * Use when you don't want to block execution
   * @param {string} category - Log category
   * @param {string} level - Log level
   * @param {string} message - Log message
   * @param {Object} context - Additional context
   */
  fireAndForget(category, level, message, context = {}) {
    this.log(category, level, message, context).catch(() => {
      // Silently ignore errors in fire-and-forget mode
    });
  }

  /**
   * Express middleware for automatic request logging
   * @returns {Function} Express middleware
   */
  expressMiddleware() {
    return async (req, res, next) => {
      const start = Date.now();
      
      res.on('finish', () => {
        const duration = Date.now() - start;
        this.fireAndForget(
          OneLog.Categories.PERFORMANCE,
          duration > 2000 ? OneLog.Levels.WARN : OneLog.Levels.INFO,
          `${req.method} ${req.path} - ${res.statusCode}`,
          {
            endpoint: req.path,
            method: req.method,
            status_code: res.statusCode,
            duration_ms: duration,
            ip_address: req.ip
          }
        );
      });

      next();
    };
  }

  /**
   * Get client IP address
   * @private
   * @returns {string|null}
   */
  _getIPAddress() {
    // In a real implementation, this would detect the actual IP
    return null;
  }
}

module.exports = OneLog;
