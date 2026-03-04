import React, { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import s from './Coins.module.css';
import { getSymbols } from '../../api/market';
import { money } from '../../utils/format';

import Header from '../../components/Header/Header';

export default function Coins() {
  const navigate = useNavigate();
  const [q, setQ] = useState('');
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchData = () => {
      getSymbols()
        .then(setRows)
        .catch(console.error); // Silently fail or log error
    };

    setLoading(true);
    fetchData(); // Initial fetch
    setLoading(false);

    const interval = setInterval(fetchData, 5000); // Poll every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const filtered = useMemo(() => {
    const v = q.trim().toLowerCase();
    if (!v) return rows;
    return rows.filter(r =>
      r.symbol.toLowerCase().includes(v) ||
      r.name.toLowerCase().includes(v)
    );
  }, [q, rows]);

  return (
    <div className={s.wrap}>
      <Header>
        <Link to="/profile" className={s.profileBtn}>
          Профиль →
        </Link>
      </Header>

      <div className={s.searchRow}>
        <input
          className={s.search}
          placeholder="Поиск по монетам"
          value={q}
          onChange={e => setQ(e.target.value)}
        />
      </div>

      <section className={s.tableWrap}>
        <table className={s.table}>
          <thead>
            <tr>
              <th>#</th>
              <th>Символ</th>
              <th>Название</th>
              <th>Цена, EDU</th>
              <th>PnL</th>
              <th>Открыть</th>
            </tr>
          </thead>
          <tbody>
            {loading && rows.length === 0 && (
              <tr><td colSpan={6} className={s.centerMuted}>Загрузка…</td></tr>
            )}
            {!loading && filtered.length === 0 && rows.length > 0 && (
              <tr><td colSpan={6} className={s.centerMuted}>Ничего не найдено</td></tr>
            )}
            {filtered.length === 0 && rows.length === 0 && !loading && (
              <tr><td colSpan={6} className={s.centerMuted}>Нет данных</td></tr>
            )}
            {filtered.map((r, i) => (
              <tr key={r.symbol}>
                <td>{i + 1}</td>
                <td className={s.symbol}>{r.symbol}</td>
                <td className={s.muted}>{r.name}</td>
                <td className={s.num}>{money(r.price)}</td>
                <td className={
                  r.pnl > 0 ? s.green : (r.pnl < 0 ? s.red : s.neutral)
                }>
                  {r.pnl > 0 ? '+' : ''}{Number(r.pnl).toFixed(2)}%
                </td>
                <td>
                  <button className={s.openBtn} onClick={() => navigate(`/coin/${r.symbol}`)}>
                    →
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </div>
  );
}
