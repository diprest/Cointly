import React, { useState } from 'react';
import s from '../Coin.module.css';

export default function PredictionWidget({ preds, onPredict }) {
  const [goal, setGoal] = useState('');
  const [horiz, setHoriz] = useState('');
  const [prob, setProb] = useState('');

  const handleSubmit = () => {
    onPredict({ goal, horiz, prob });
    setGoal(''); setHoriz(''); setProb('');
  };

  const predsView = preds.slice(0, 15);

  return (
    <>
      <section className={`${s.card} ${s.predCard}`}>
        <div className={s.cardTitle}>Последние прогнозы по монете</div>
        <div className={s.predList}>
          {predsView.map(p => (
            <div key={p.id} className={s.predRow}>
              <div className={s.author}>{p.author}</div>
              <div className={s.predText}>
                Цель {p.target.toLocaleString('ru-RU')} EDU • Горизонт {p.horizonDays} дн. • Вероятность {(p.prob * 100).toFixed(0)}%
              </div>
              <div className={p.result >= 0 ? s.green : s.red}>
                {p.result >= 0 ? '+' : ''}{p.result}%
              </div>
            </div>
          ))}
        </div>
      </section>

      <section className={`${s.card} ${s.formCard}`}>
        <div className={s.cardTitle}>Прогноз</div>
        <div className={s.formStack}>
          <div className={s.formRow}>
            <label>Цель, EDU</label>
            <input className={s.input} value={goal} onChange={e => setGoal(e.target.value)} placeholder="напр., 75000" />
          </div>
          <div className={s.formRow}>
            <label>Горизонт</label>
            <input className={s.input} value={horiz} onChange={e => setHoriz(e.target.value)} placeholder="дней" />
          </div>
          <div className={s.formRow}>
            <label>Вероятность</label>
            <input className={s.input} value={prob} onChange={e => setProb(e.target.value)} placeholder="%" />
          </div>
        </div>
        <div className={s.actionsRight}>
          <button className={s.createBtn} onClick={handleSubmit}>Создать прогноз</button>
        </div>
      </section>
    </>
  );
}
