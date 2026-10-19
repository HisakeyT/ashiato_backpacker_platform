export const formatPeriod = (start: string | null, end: string | null): string => {
  if (!start) return "日付未定";
  if (!end || end === start) return start;
  return `${start} 〜 ${end}`;
};
