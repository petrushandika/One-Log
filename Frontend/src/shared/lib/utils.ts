/** Maps raw log category enum values to human-readable labels. */
export const CATEGORY_LABELS: Record<string, string> = {
  SYSTEM_ERROR:  'System Error',
  AUTH_EVENT:    'Auth Event',
  USER_ACTIVITY: 'User Activity',
  SECURITY:      'Security',
  PERFORMANCE:   'Performance',
  AUDIT_TRAIL:   'Audit Trail',
};

export function categoryLabel(cat: string): string {
  return CATEGORY_LABELS[cat] ?? cat;
}
