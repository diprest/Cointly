import React, { useMemo, useState, useEffect } from "react";
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import s from "./ProfileWallet.module.css";
import Header from '../../components/Header/Header';
import ProfileNavigation from '../../components/ProfileNavigation/ProfileNavigation';
import { getPortfolio } from '../../api/portfolio';
import { getSymbols } from '../../api/market';

const ProfileWallet = ({ userId }) => {
  const navigate = useNavigate();
  const [q, setQ] = useState("");
  const [portfolio, setPortfolio] = useState([]);
  const [marketData, setMarketData] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let interval;
    const loadData = async () => {
      if (!userId) return;
      try {
        const [portData, markData] = await Promise.all([
          getPortfolio(userId),
          getSymbols()
        ]);
        setPortfolio(portData || []);
        setMarketData(markData || []);
        setLoading(false);
      } catch (e) {
        console.error('[ProfileWallet] Error:', e);
        setLoading(false);
      }
    };

    loadData();
    interval = setInterval(loadData, 10000);
    return () => clearInterval(interval);
  }, [userId]);

  const stats = useMemo(() => {
    let portfolioValue = 0;
    let balance = 0;
    const assets = [];

    portfolio.forEach(item => {
      const amount = parseFloat(item.amount || 0);
      const totalCost = parseFloat(item.total_cost || 0);

      if (item.asset === 'USDT') {
        balance += amount;
      } else {
        let price = 0;
        const ticker = marketData.find(m => m.symbol === item.asset + 'USDT');
        price = ticker ? ticker.price : 0;

        const currentValue = amount * price;
        portfolioValue += currentValue;

        const coinPnlAbs = currentValue - totalCost;
        const coinPnlPct = totalCost > 0 ? (coinPnlAbs / totalCost) * 100 : 0;

        if (amount > 0) {
          assets.push({
            id: item.asset,
            symbol: item.asset,
            qty: amount,
            spent: totalCost,
            price: price,
            value: currentValue,
            pnlAbs: coinPnlAbs,
            pnlPct: coinPnlPct
          });
        }
      }
    });

    const totalEquity = portfolioValue + balance;
    const initialCapital = 10000;
    const totalPnlAbs = totalEquity - initialCapital;
    const totalPnlPct = (totalPnlAbs / initialCapital) * 100;

    return { portfolioValue, balance, totalPnlAbs, totalPnlPct, assets };
  }, [portfolio, marketData]);

  const rows = useMemo(() => {
    const v = q.trim().toLowerCase();
    let data = stats.assets;
    if (v) {
      data = data.filter((r) => r.symbol.toLowerCase().includes(v));
    }
    return data.sort((a, b) => b.value - a.value);
  }, [q, stats.assets]);

  return (
    <div className={s.wrap}>
      <Header>
        <ProfileNavigation />
      </Header>

      <main className={s.page}>
        <div className={s.grid}>
          <section className={s.topRow}>
            <div className={s.cardStat}>
              <div className={s.caption}>Стоимость портфеля</div>
              <div className={s.value}>
                {loading ? '...' : `$${stats.portfolioValue.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}
              </div>
              <div className={s.subValue} style={{ fontSize: '14px', color: '#888', marginTop: '4px' }}>
                Баланс: {loading ? '...' : `$${stats.balance.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}
              </div>
            </div>

            <div className={s.cardStatWide}>
              <div className={s.caption}>Общий PnL (от $10,000)</div>
              <div className={stats.totalPnlAbs >= 0 ? s.valueGreen : s.valueRed}>
                {loading ? '...' : `${stats.totalPnlPct >= 0 ? '+' : ''}${stats.totalPnlPct.toFixed(2)}%`}
                <span style={{ fontSize: '0.6em', marginLeft: '10px', opacity: 0.8 }}>
                  ({loading ? '...' : `${stats.totalPnlAbs >= 0 ? '+' : ''}$${Math.abs(stats.totalPnlAbs).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`})
                </span>
              </div>
            </div>
          </section>

          <div className={s.searchWrap}>
            <input
              className={s.search}
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="Поиск монеты..."
            />
          </div>

          <section className={s.tableCard}>
            <div className={`${s.tr} ${s.head}`}>
              <div>Монета</div>
              <div className={s.right}>Кол-во</div>
              <div className={s.right}>Цена</div>
              <div className={s.right}>Стоимость</div>
              <div className={s.right}>Ср. Цена</div>
              <div className={s.right}>PnL</div>
              <div className={s.right}>PnL%</div>
              <div className={s.center}>Открыть</div>
            </div>

            <div className={s.body}>
              {rows.length > 0 ? (
                rows.map((r) => (
                  <div className={s.tr} key={r.id}>
                    <div className={s.bold}>{r.symbol}</div>
                    <div className={`${s.mono} ${s.right}`}>{r.qty.toLocaleString("en-US", { maximumFractionDigits: 6 })}</div>
                    <div className={`${s.mono} ${s.right}`}>${r.price.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</div>
                    <div className={`${s.mono} ${s.right}`}>${r.value.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</div>
                    <div className={`${s.mono} ${s.right}`}>
                      ${(r.qty > 0 ? r.spent / r.qty : 0).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </div>
                    <div className={`${r.pnlAbs >= 0 ? s.green : s.red} ${s.right}`}>
                      {r.pnlAbs >= 0 ? '+' : ''}${Math.abs(r.pnlAbs).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </div>
                    <div className={`${r.pnlPct >= 0 ? s.green : s.red} ${s.right}`}>
                      {r.pnlPct >= 0 ? '+' : ''}{r.pnlPct.toFixed(2)}%
                    </div>
                    <div className={s.center}>
                      <button
                        className={s.openBtn}
                        onClick={() => navigate(`/coin/${r.symbol}USDT`)}
                        aria-label="Открыть"
                      >
                        ›
                      </button>
                    </div>
                  </div>
                ))
              ) : (
                <div className={s.empty}>
                  {loading ? 'Загрузка...' : 'Нет активов'}
                </div>
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

export default connect(mapStateToProps)(ProfileWallet);
