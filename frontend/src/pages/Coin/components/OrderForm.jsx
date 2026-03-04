import React, { useState, useMemo } from 'react';
import { toast } from 'react-toastify';
import s from '../Coin.module.css';

export default function OrderForm({ lastPrice, onOrder }) {
  const [side, setSide] = useState('buy');
  const [qty, setQty] = useState('0.100');
  const [orderType, setOrderType] = useState('market');
  const [price, setPrice] = useState('');

  const currentPrice = orderType === 'market' ? lastPrice : (parseFloat(price) || 0);
  const total = (parseFloat(qty) || 0) * currentPrice;

  return (
    <section className={`${s.card} ${s.orderCard}`}>
      <div className={s.cardTitle}>Ордер</div>

      <div className={s.cardContent}>
        <div className={s.switchRow}>
          <button
            className={`${s.switchBtn} ${side === 'buy' ? s.switchActive : ''}`}
            onClick={() => setSide('buy')}
          >
            Купить
          </button>
          <button
            className={`${s.switchBtn} ${side === 'sell' ? s.switchActive : ''}`}
            onClick={() => setSide('sell')}
          >
            Продать
          </button>
        </div>

        <div className={s.formRow}>
          <label>Количество</label>
          <input
            type="number"
            value={qty}
            onChange={e => setQty(e.target.value)}
            className={s.input}
          />
        </div>

        <div className={s.formRow}>
          <label>Цена</label>
          <div className={s.splitInput}>
            <select
              className={s.select}
              value={orderType}
              onChange={e => setOrderType(e.target.value)}
            >
              <option value="market">Рыночная</option>
              <option value="limit">Ордер</option>
            </select>
            <input
              type="number"
              value={orderType === 'market' ? 'Рыночная' : price}
              onChange={e => setPrice(e.target.value)}
              className={s.input}
              placeholder={orderType === 'market' ? 'Рыночная' : 'Цена'}
              disabled={orderType === 'market'}
              readOnly={orderType === 'market'}
            />
          </div>
        </div>

        <div className={s.totalRow}>
          Итого <span>≈ {total.toLocaleString('ru-RU', { maximumFractionDigits: 2 })} USDT</span>
        </div>

        <div className={s.actions}>
          <button
            className={side === 'buy' ? s.btnBuy : s.btnSell}
            onClick={() => {
              if (!qty || parseFloat(qty) <= 0) {
                toast.error("Введите корректное количество");
                return;
              }
              onOrder(side, qty, orderType === 'market' ? null : price);
            }}
          >
            {side === 'buy' ? 'Купить' : 'Продать'}
          </button>
          <button className={s.btnGhost} onClick={() => setQty('0.100')}>
            Отмена
          </button>
        </div>
      </div>
    </section>
  );
}
