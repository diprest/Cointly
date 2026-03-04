import React, { useState, useEffect } from 'react';
import s from '../Coin.module.css';

const TimeLeft = ({ expiresAt }) => {
  const [left, setLeft] = useState('');

  useEffect(() => {
    const update = () => {
      const now = new Date();
      const end = new Date(expiresAt);
      const diff = Math.floor((end - now) / 1000);
      if (diff <= 0) setLeft('0s');
      else setLeft(`${diff}s`);
    };
    update();
    const interval = setInterval(update, 1000);
    return () => clearInterval(interval);
  }, [expiresAt]);

  return <span>{left}</span>;
};

export default function ActiveBets({ bets }) {
  return (
    <div className={s.card}>
      <div className={s.cardTitle}>Мои активные ставки</div>

      <div className={`${s.tableHeader} ${s.betsGrid}`}>
        <div>Направление</div>
        <div>Цена старта</div>
        <div>Сумма</div>
        <div className={s.colRight}>Таймер</div>
      </div>

      <div className={s.ordersList}>
        {bets.length === 0 && (
          <div className={s.emptyState}>
            <div className={s.emptyIcon}>🎲</div>
            <div>Нет активных ставок</div>
          </div>
        )}

        {bets.map((bet) => (
          <div key={bet.id} className={`${s.orderItem} ${s.betsGrid}`}>
            <div className={bet.direction === 'UP' ? s.green : s.red}>
              {bet.direction === 'UP' ? 'Вверх' : 'Вниз'}
            </div>
            <div className={s.mono}>{bet.opened_price}</div>
            <div className={s.mono}>{bet.stake_amount} EDU</div>
            <div className={`${s.mono} ${s.colRight}`}>
              <TimeLeft expiresAt={bet.expires_at} />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
