import React, { useEffect, useRef, useMemo } from 'react';
import s from '../Coin.module.css';
import { drawLine } from '../../../utils/chart';

const TICK_MS = 1200;

export default function CoinChart({ series }) {
  const canvasRef = useRef(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || !series.length) return;

    const host = canvas.parentElement;
    const rect = host.getBoundingClientRect();

    const left = 68, right = 14, top = 12, bottom = 36;

    const w = rect.width - left - right;
    const h = rect.height - top - bottom;
    canvas.style.position = 'absolute';
    canvas.style.left = `${left}px`;
    canvas.style.top = `${top}px`;
    canvas.style.width = `${w}px`;
    canvas.style.height = `${h}px`;

    drawLine(canvas, series, '#e6dd9b', { top: 0, bottom: 0, left: 0, right: 0 });
  }, [series]);

  const yTicks = useMemo(() => {
    if (!series.length) return [];
    const ys = series.map(p => p.y);
    const min = Math.min(...ys), max = Math.max(...ys);
    const steps = 5;
    return new Array(steps + 1).fill(0).map((_, i) => {
      const t = i / steps;
      return { label: (min + (max - min) * (1 - t)).toFixed(0), pos: t };
    });
  }, [series]);

  const xTicks = useMemo(() => {
    if (!series.length) return [];
    const steps = 6;
    const totalMs = (series.length - 1) * TICK_MS;
    const start = Date.now() - totalMs;
    let last = '';
    return new Array(steps + 1).fill(0).map((_, i) => {
      const t = i / steps;
      const ts = start + t * totalMs;
      const label = new Date(ts).toLocaleTimeString('ru-RU', { minute: '2-digit', second: '2-digit' });
      const out = label === last ? '\u00A0' : label;
      last = label;
      return { label: out, pos: t };
    });
  }, [series]);

  return (
    <section className={`${s.card} ${s.chartCard}`}>
      <div className={s.cardTitle}>Цена, EDU</div>
      <div className={s.chartInner}>
        <div className={s.gridH}>
          {yTicks.map((t, i) => (
            <div key={i} className={s.gridLine} style={{ top: `calc(${t.pos * 100}% )` }} />
          ))}
        </div>
        <div className={s.gridV}>
          {xTicks.map((t, i) => (
            <div key={i} className={s.gridVLine} style={{ left: `calc(${t.pos * 100}% )` }} />
          ))}
        </div>

        <div className={s.axisY}>
          {yTicks.map((t, i) => (
            <div key={i}>{t.label}</div>
          ))}
        </div>
        <div className={s.axisX}>
          {xTicks.map((t, i) => (
            <div key={i}>{t.label}</div>
          ))}
        </div>

        <canvas ref={canvasRef} className={s.chartCanvas} />
      </div>
    </section>
  );
}
