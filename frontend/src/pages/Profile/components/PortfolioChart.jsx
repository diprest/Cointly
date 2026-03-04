import React, { useEffect, useRef } from 'react';
import s from '../Profile.module.css';
import { drawDonut, drawLine } from '../../../utils/chart';

export function PortfolioDonut({ allocations, valueLabel }) {
  const donutRef = useRef(null);

  useEffect(() => {
    const donut = donutRef.current;
    if (donut) {
      drawDonut(donut, allocations, valueLabel);
    }
  }, [allocations, valueLabel]);

  return (
    <div className={s.donutWrap}>
      <div className={s.donutBox}>
        <canvas ref={donutRef} className={s.donutCanvas} />
      </div>
      <ul className={s.legend}>
        {allocations.map(i => (
          <li key={i.symbol}>
            <span className={s.dot} style={{ background: i.color }} />
            {i.name}
          </li>
        ))}
      </ul>
    </div>
  );
}

export function PortfolioHistory({ series, period }) {
  const lineRef = useRef(null);

  useEffect(() => {
    const line = lineRef.current;
    if (line) {
      drawLine(line, series);
    }
  }, [series, period]);

  return (
    <div className={s.lineWrap}>
      <canvas ref={lineRef} className={s.lineCanvas} />
    </div>
  );
}
