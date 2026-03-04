export const money = (v) =>
  (typeof v === 'number' ? v : +v).toLocaleString('ru-RU', { minimumFractionDigits: 2, maximumFractionDigits: 2 });

export const pct = (v) => {
  if (v === 0) return '0.0%';
  const sign = v > 0 ? '+' : '';
  return `${sign}${(v * 100).toFixed(1)}%`;
};
