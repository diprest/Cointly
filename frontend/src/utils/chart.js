export function drawDonut(canvas, items, centerLabel) {
  const ctx = canvas.getContext('2d');

  const W = canvas.clientWidth;
  const H = canvas.clientHeight;
  canvas.width = W;
  canvas.height = H;
  ctx.clearRect(0, 0, W, H);

  const r = Math.min(W, H) * 0.36;
  const cx = W * 0.5, cy = H * 0.5;
  let start = -Math.PI / 2;

  items.forEach(it => {
    const ang = it.share * Math.PI * 2;
    ctx.beginPath();
    ctx.moveTo(cx, cy);
    ctx.arc(cx, cy, r, start, start + ang, false);
    ctx.fillStyle = it.color;
    ctx.fill();
    start += ang;
  });

  ctx.globalCompositeOperation = 'destination-out';
  ctx.beginPath();
  ctx.arc(cx, cy, r * 0.58, 0, Math.PI * 2);
  ctx.fill();
  ctx.globalCompositeOperation = 'source-over';

  ctx.fillStyle = '#e7e7ea';
  ctx.font = '700 16px system-ui';
  ctx.textAlign = 'center';
  ctx.fillText('Портфель', cx, cy - 6);
  ctx.font = '800 18px system-ui';
  ctx.fillText(centerLabel, cx, cy + 16);
}

export function drawLine(canvas, points, color = '#e6dd9b', padding = { top: 8, bottom: 8, left: 8, right: 8 }) {
  const ctx = canvas.getContext('2d');

  const W = canvas.clientWidth;
  const H = canvas.clientHeight;
  canvas.width = W;
  canvas.height = H;
  ctx.clearRect(0, 0, W, H);

  ctx.fillStyle = '#2a2f36'; ctx.fillRect(0, 0, W, H);

  if (!points.length) return;
  const ys = points.map(p => p.y !== undefined ? p.y : p);
  const xs = points.map((p, i) => p.x !== undefined ? p.x : i);

  const minY = Math.min(...ys), maxY = Math.max(...ys);
  const { top: padT, bottom: padB, left: padL, right: padR } = padding;
  const innerW = W - padL - padR, innerH = H - padT - padB;

  const X = (val) => padL + (val - xs[0]) / (xs[xs.length - 1] - xs[0] || 1) * innerW;
  const Y = (val) => padT + (1 - ((val - minY) / ((maxY - minY) || 1))) * innerH;

  ctx.lineWidth = 3; ctx.strokeStyle = color;
  ctx.beginPath(); ctx.moveTo(X(xs[0]), Y(ys[0]));
  for (let i = 1; i < points.length; i++) ctx.lineTo(X(xs[i]), Y(ys[i]));
  ctx.stroke();

  ctx.lineTo(X(xs[points.length - 1]), H - padB);
  ctx.lineTo(X(xs[0]), H - padB);
  ctx.closePath(); ctx.fillStyle = 'rgba(230,221,155,.18)'; ctx.fill();
}
