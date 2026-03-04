import React, { useState, useEffect } from "react";
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import s from "./ProfileHistory.module.css";
import Header from '../../components/Header/Header';
import ProfileNavigation from '../../components/ProfileNavigation/ProfileNavigation';
import { getOrders } from '../../api/trading';
import { getUserBets } from '../../api/bets';

const Pill = ({ kind, children }) => (
  <span className={`${s.pill} ${s[kind]}`}>{children}</span>
);

const Countdown = ({ expiresAt }) => {
  const [timeLeft, setTimeLeft] = useState(0);

  useEffect(() => {
    const update = () => {
      const now = new Date().getTime();
      const end = new Date(expiresAt).getTime();
      const diff = Math.max(0, Math.floor((end - now) / 1000));
      setTimeLeft(diff);
    };
    update();
    const interval = setInterval(update, 1000);
    return () => clearInterval(interval);
  }, [expiresAt]);

  return <span>{timeLeft}s</span>;
};

const ProfileHistory = ({ userId }) => {
  const navigate = useNavigate();
  const [orders, setOrders] = useState([]);
  const [bets, setBets] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let interval;
    const loadData = async () => {
      if (!userId) return;
      try {
        const [ordersData, betsData] = await Promise.all([
          getOrders(userId),
          getUserBets(userId)
        ]);
        setOrders(ordersData || []);
        setBets(betsData || []);
        setLoading(false);
      } catch (e) {
        console.error('[ProfileHistory] Error:', e);
        setLoading(false);
      }
    };

    loadData();
    interval = setInterval(loadData, 5000);
    return () => clearInterval(interval);
  }, [userId]);

  const activeOrders = orders.filter(o => o.status === 'NEW' || o.status === 'PARTIAL');
  const historyOrders = orders.filter(o => o.status !== 'NEW' && o.status !== 'PARTIAL').sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

  const activeBets = bets.filter(b => b.status === 'OPEN');
  const historyBets = bets.filter(b => b.status !== 'OPEN').sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

  const statusPill = (st) => {
    if (st === "FILLED") return <Pill kind="ok">Исполнено</Pill>;
    if (st === "PARTIAL") return <Pill kind="warn">Частично</Pill>;
    if (st === "CANCELED") return <Pill kind="bad">Отменено</Pill>;
    if (st === "NEW") return <Pill kind="warn">Активен</Pill>;
    return <Pill kind="bad">{st}</Pill>;
  };

  const formatDate = (dateStr) => {
    if (!dateStr) return '-';
    const d = new Date(dateStr);
    const pad = (n) => n.toString().padStart(2, '0');
    return `${pad(d.getDate())}.${pad(d.getMonth() + 1)}.${d.getFullYear()}, ${pad(d.getHours())}:${pad(d.getMinutes())}`;
  };

  return (
    <div className={s.wrap}>
      <Header>
        <ProfileNavigation />
      </Header>

      <main className={s.page}>
        <div className={s.grid}>

          <section className={s.card}>
            <div className={s.cardTitle}>Активные ордера</div>
            <div className={s.activeList}>
              {activeOrders.length > 0 ? (
                activeOrders.map((o) => (
                  <div key={o.id} className={s.activeItem}>
                    <div className={s.activeHead}>
                      <div
                        className={s.activeCoin}
                        onClick={() => navigate(`/coin/${o.symbol}`)}
                        style={{ cursor: 'pointer', textDecoration: 'underline' }}
                      >
                        {o.symbol}
                      </div>
                    </div>
                    <div className={s.activeText}>
                      {o.side} {o.amount} по {o.price}
                    </div>
                  </div>
                ))
              ) : (
                <div className={s.empty}>Нет активных ордеров</div>
              )}
            </div>
          </section>

          <section className={s.card}>
            <div className={s.cardTitle}>История сделок</div>

            <div className={`${s.tr} ${s.head} ${s.trTrades}`}>
              <div>Дата</div>
              <div>Монета</div>
              <div>Тип</div>
              <div>Операция</div>
              <div>Кол-во</div>
              <div>Цена</div>
              <div>Статус</div>
            </div>

            <div className={s.tableBody}>
              {historyOrders.length > 0 ? (
                historyOrders.map((r) => (
                  <div key={r.id} className={`${s.tr} ${s.trTrades}`}>
                    <div className={s.mono}>{formatDate(r.created_at)}</div>
                    <div
                      className={s.bold}
                      onClick={() => navigate(`/coin/${r.symbol}`)}
                      style={{ cursor: 'pointer', textDecoration: 'underline' }}
                    >
                      {r.symbol}
                    </div>
                    <div className={s.mono}>{r.type}</div>
                    <div className={r.side === "SELL" ? s.red : s.green}>
                      {r.side === "SELL" ? "Продажа" : "Покупка"}
                    </div>
                    <div className={s.mono}>{r.amount}</div>
                    <div className={s.mono}>{r.price}</div>
                    <div>{statusPill(r.status)}</div>
                  </div>
                ))
              ) : (
                <div className={s.empty}>История пуста</div>
              )}
            </div>
          </section>

          <section className={s.card}>
            <div className={s.cardTitle}>Активные ставки</div>
            <div className={s.activeList}>
              {activeBets.length > 0 ? (
                activeBets.map((a) => (
                  <div key={a.id} className={s.activeItem}>
                    <div className={s.activeHead}>
                      <div
                        className={s.activeCoin}
                        onClick={() => navigate(`/coin/${a.symbol}`)}
                        style={{ cursor: 'pointer', textDecoration: 'underline' }}
                      >
                        {a.symbol}
                      </div>
                      <div className={s.mono}>
                        <Countdown expiresAt={a.expires_at} />
                      </div>
                    </div>
                    <div className={s.activeText}>
                      {a.stake_amount} • <span className={a.direction === 'UP' ? s.green : s.red}>{a.direction === 'UP' ? 'Вверх' : 'Вниз'}</span>
                    </div>
                  </div>
                ))
              ) : (
                <div className={s.empty}>Нет активных ставок</div>
              )}
            </div>
          </section>

          <section className={s.card}>
            <div className={s.cardTitle}>История ставок</div>

            <div className={`${s.tr} ${s.head} ${s.trPreds}`}>
              <div>Дата</div>
              <div>Монета</div>
              <div>Сумма</div>
              <div>Направление</div>
              <div>Результат</div>
            </div>

            <div className={s.tableBody}>
              {historyBets.length > 0 ? (
                historyBets.map((r) => (
                  <div key={r.id} className={`${s.tr} ${s.trPreds}`}>
                    <div className={s.mono}>{formatDate(r.created_at)}</div>
                    <div
                      className={s.bold}
                      onClick={() => navigate(`/coin/${r.symbol}`)}
                      style={{ cursor: 'pointer', textDecoration: 'underline' }}
                    >
                      {r.symbol}
                    </div>
                    <div className={s.mono}>{r.stake_amount}</div>
                    <div className={r.direction === 'UP' ? s.green : s.red}>{r.direction}</div>
                    <div className={r.status === "WON" ? s.green : s.red}>
                      {r.status === "WON" ? `+${(r.stake_amount * 0.8).toFixed(2)}` : `-${r.stake_amount}`}
                    </div>
                  </div>
                ))
              ) : (
                <div className={s.empty}>История ставок пуста</div>
              )}
            </div>
          </section>

        </div>
      </main>
    </div>
  );
};

const mapStateToProps = (state) => ({
  userId: state.user.userId
});

export default connect(mapStateToProps)(ProfileHistory);
