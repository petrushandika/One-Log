"""
One-Log (ULAM) Official Python SDK

Example usage:
    from onelog import OneLog

    client = OneLog(
        api_key="ulam_live_xxxxxxxx",
        endpoint="https://api.ulam.your-domain.com"
    )

    client.log_info("SYSTEM_ERROR", "Database connection established")
"""

import os
import json
import threading
import requests
from enum import Enum
from typing import Dict, Any, Optional
from datetime import datetime


class LogLevel(str, Enum):
    """Log severity levels"""
    CRITICAL = "CRITICAL"
    ERROR = "ERROR"
    WARN = "WARN"
    INFO = "INFO"
    DEBUG = "DEBUG"


class Category(str, Enum):
    """Log categories"""
    SYSTEM_ERROR = "SYSTEM_ERROR"
    USER_ACTIVITY = "USER_ACTIVITY"
    AUTH_EVENT = "AUTH_EVENT"
    PERFORMANCE = "PERFORMANCE"
    SECURITY = "SECURITY"
    AUDIT_TRAIL = "AUDIT_TRAIL"


class OneLog:
    """
    One-Log (ULAM) client for sending logs to the platform.
    
    Args:
        api_key: Your ULAM API key (or set ULAM_API_KEY env var)
        endpoint: API endpoint URL (or set ULAM_ENDPOINT env var)
        source_id: Source identifier (or set ULAM_SOURCE_ID env var)
        timeout: Request timeout in seconds (default: 5)
    """

    def __init__(
        self,
        api_key: Optional[str] = None,
        endpoint: Optional[str] = None,
        source_id: Optional[str] = None,
        timeout: float = 5.0
    ):
        self.api_key = api_key or os.getenv("ULAM_API_KEY")
        self.endpoint = endpoint or os.getenv("ULAM_ENDPOINT") or "http://localhost:8080/api/ingest"
        self.source_id = source_id or os.getenv("ULAM_SOURCE_ID")
        self.timeout = timeout

        if not self.api_key:
            raise ValueError("API key is required. Set ULAM_API_KEY environment variable or pass api_key parameter.")

        self.session = requests.Session()
        self.session.headers.update({
            "Content-Type": "application/json",
            "X-API-Key": self.api_key
        })

    def log(
        self,
        category: Category,
        level: LogLevel,
        message: str,
        context: Optional[Dict[str, Any]] = None,
        stack_trace: Optional[str] = None
    ) -> None:
        """
        Send a log entry to One-Log.

        Args:
            category: Log category
            level: Log severity level
            message: Log message
            context: Additional context data
            stack_trace: Stack trace (optional)
        """
        entry = {
            "category": category.value,
            "level": level.value,
            "message": message,
            "context": context or {},
            "stack_trace": stack_trace,
            "ip_address": self._get_ip_address()
        }

        try:
            response = self.session.post(
                self.endpoint,
                json=entry,
                timeout=self.timeout
            )
            response.raise_for_status()
        except requests.RequestException as e:
            raise Exception(f"Failed to send log: {e}")

    def log_with_error(
        self,
        category: Category,
        level: LogLevel,
        message: str,
        error: Exception,
        context: Optional[Dict[str, Any]] = None
    ) -> None:
        """
        Log an error with its details.

        Args:
            category: Log category
            level: Log severity level
            message: Log message
            error: Exception object
            context: Additional context data
        """
        import traceback
        stack_trace = "".join(traceback.format_exception(type(error), error, error.__traceback__))
        self.log(category, level, message, context, stack_trace)

    def log_critical(self, category: Category, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log a critical level message"""
        self.log(category, LogLevel.CRITICAL, message, context)

    def log_error(self, category: Category, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log an error level message"""
        self.log(category, LogLevel.ERROR, message, context)

    def log_warn(self, category: Category, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log a warning level message"""
        self.log(category, LogLevel.WARN, message, context)

    def log_info(self, category: Category, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log an info level message"""
        self.log(category, LogLevel.INFO, message, context)

    def log_debug(self, category: Category, message: str, context: Optional[Dict[str, Any]] = None) -> None:
        """Log a debug level message"""
        self.log(category, LogLevel.DEBUG, message, context)

    def log_performance(
        self,
        endpoint: str,
        duration_ms: int,
        method: str = "GET",
        status_code: int = 200
    ) -> None:
        """
        Log performance metrics.

        Args:
            endpoint: API endpoint path
            duration_ms: Request duration in milliseconds
            method: HTTP method
            status_code: HTTP status code
        """
        self.log(
            Category.PERFORMANCE,
            LogLevel.INFO,
            "API request completed",
            {
                "endpoint": endpoint,
                "duration_ms": duration_ms,
                "method": method,
                "status_code": status_code
            }
        )

    def log_auth_event(
        self,
        event_type: str,
        auth_method: str,
        user_id: str,
        success: bool = True
    ) -> None:
        """
        Log an authentication event.

        Args:
            event_type: Event type (e.g., 'login_success', 'login_failed')
            auth_method: Authentication method
            user_id: User ID
            success: Whether authentication succeeded
        """
        level = LogLevel.INFO if success else LogLevel.WARN
        self.log(
            Category.AUTH_EVENT,
            level,
            "Authentication event",
            {
                "event_type": event_type,
                "auth_method": auth_method,
                "user_id": user_id,
                "success": success
            }
        )

    def log_audit(
        self,
        action: str,
        actor_id: str,
        resource_type: str,
        resource_id: str,
        before: Any = None,
        after: Any = None
    ) -> None:
        """
        Log an audit trail event.

        Args:
            action: Action performed
            actor_id: Actor/user ID
            resource_type: Type of resource
            resource_id: Resource ID
            before: State before change
            after: State after change
        """
        self.log(
            Category.AUDIT_TRAIL,
            LogLevel.INFO,
            "Audit event",
            {
                "action": action,
                "actor_id": actor_id,
                "resource_type": resource_type,
                "resource_id": resource_id,
                "before": before,
                "after": after
            }
        )

    def fire_and_forget(
        self,
        category: Category,
        level: LogLevel,
        message: str,
        context: Optional[Dict[str, Any]] = None
    ) -> None:
        """
        Fire-and-forget logging (async, non-blocking).
        Use when you don't want to block execution.
        """
        def send_log():
            try:
                self.log(category, level, message, context)
            except Exception:
                # Silently ignore errors in fire-and-forget mode
                pass

        thread = threading.Thread(target=send_log)
        thread.daemon = True
        thread.start()

    def _get_ip_address(self) -> Optional[str]:
        """Get client IP address"""
        return None

    def close(self) -> None:
        """Close the HTTP session"""
        self.session.close()

    def __enter__(self):
        """Context manager entry"""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit"""
        self.close()
