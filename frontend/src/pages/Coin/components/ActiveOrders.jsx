import React from 'react';
import s from '../Coin.module.css';

export default function ActiveOrders({ orders, onCancel }) {
  return (
    <div className={s.card}>
      <div className={s.cardTitle}>Активные ордера</div>

      <div className={`${s.tableHeader} ${s.ordersGrid}`}>
        <div>Тип</div>
        <div>Цена</div>
        <div>Кол-во</div>
        <div className={s.colRight}>Действие</div>
      </div>

      <div className={s.ordersList}>
        {orders.length === 0 && (
          <div className={s.emptyState}>
            <div className={s.emptyIcon}>📭</div>
            <div>Нет активных ордеров</div>
          </div>
        )}

        {orders.map((o) => (
          <div key={o.id} className={`${s.orderItem} ${s.ordersGrid}`}>
            <div className={o.side === 'BUY' ? s.green : s.red}>
              {o.side === 'BUY' ? 'Купить' : 'Продать'}
            </div>
            <div className={s.mono}>{o.price}</div>
            <div className={s.mono}>{o.amount} {o.symbol.replace('USDT', '')}</div>
            <button className={`${s.btnSmallGhost} ${s.colRight}`} onClick={() => onCancel(o.id)}>
              Отменить
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
