import React, { useState } from 'react';
import { toast } from 'react-toastify';
import s from '../Coin.module.css';

export default function BettingWidget({ onBet }) {
  const [amount, setAmount] = useState('');
  const [timeframe, setTimeframe] = useState('1m'); // 1m, 5m, 1h

  const handleBet = (direction) => {
    if (!amount) {
      toast.error("Введите сумму ставки");
      return;
    }
    const durationMap = { '1m': 60, '5m': 300, '1h': 3600 };
    const duration = durationMap[timeframe] || 60;
    onBet(direction === 'Up' ? 'up' : 'down', amount, duration);
  };

  return (
    <div className={`${s.card} ${s.orderCard}`}>
      <div className={s.cardTitle}>Сделать ставку</div>

      <div className={s.cardContent}>
        <div className={s.formRow}>
          <label>Сумма (EDU)</label>
          <input
            type="number"
            value={amount}
            onChange={e => setAmount(e.target.value)}
            placeholder="100"
            className={s.input}
          />
        </div>

        <div className={s.formRow}>
          <label>Время</label>
          <div className={s.chips}>
            {['1m', '5m', '1h'].map(tf => (
              <button
                key={tf}
                className={`${s.chip} ${timeframe === tf ? s.chipActive : ''}`}
                onClick={() => setTimeframe(tf)}
              >
                {tf}
              </button>
            ))}
          </div>
        </div>

        <div className={s.betButtons}>
          <button className={s.btnUp} onClick={() => handleBet('Up')}>
            Вверх ↗
          </button>
          <button className={s.btnDown} onClick={() => handleBet('Down')}>
            Вниз ↘
          </button>
        </div>
      </div>
    </div>
  );
}
