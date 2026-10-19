export const formatPeriod = (start: string, end: string | null): string => {
  if (!end || end === start) return start;
  return `${start} 〜 ${end}`;
};
