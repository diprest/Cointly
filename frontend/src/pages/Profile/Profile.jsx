import React from 'react';
import { connect } from 'react-redux';
import s from './Profile.module.css';
import Header from '../../components/Header/Header';
import ProfileNavigation from '../../components/ProfileNavigation/ProfileNavigation';
import { getPortfolio } from '../../api/portfolio';
import { getSymbols } from '../../api/market';

class Profile extends React.Component {
  state = {
    portfolio: [],
    marketData: [],
    loading: true
  };

  componentDidMount() {
    this.loadData();
    this.interval = setInterval(() => this.loadData(), 10000);
  }

  componentWillUnmount() {
    if (this.interval) clearInterval(this.interval);
  }

  loadData = async () => {
    const { userId } = this.props;
    if (!userId) return;

    try {
      const [portData, markData] = await Promise.all([
        getPortfolio(userId),
        getSymbols()
      ]);

      this.setState({
        portfolio: portData || [],
        marketData: markData || [],
        loading: false
      });
    } catch (e) {
      console.error('[Profile] Error:', e);
      this.setState({ loading: false });
    }
  };

  calculateStats = () => {
    const { portfolio, marketData } = this.state;

    let portfolioValue = 0;
    let balance = 0;
    const coinStats = [];

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
          coinStats.push({
            symbol: item.asset,
            title: item.asset,
            pnlAbs: coinPnlAbs,
            pnlPct: coinPnlPct,
            value: currentValue
          });
        }
      }
    });

    const totalEquity = portfolioValue + balance;
    const initialCapital = 10000;
    const pnlAbs = totalEquity - initialCapital;
    const pnlPct = (pnlAbs / initialCapital) * 100;
    const bestCoins = coinStats.sort((a, b) => b.pnlPct - a.pnlPct).slice(0, 3);

    return { portfolioValue, balance, pnlAbs, pnlPct, bestCoins };
  };

  render() {
    const { loading } = this.state;
    const stats = this.calculateStats();
    const { portfolioValue, balance, pnlAbs, pnlPct, bestCoins } = stats;

    return (
      <div className={s.wrap}>
        <Header>
          <ProfileNavigation />
        </Header>

        <div className={s.main}>
          <main className={s.grid}>
            <section className={s.cardStat}>
              <div className={s.caption}>Стоимость портфеля</div>
              <div className={s.value}>
                {loading ? '...' : `$${portfolioValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}
              </div>
              <div className={s.subValue}>
                Баланс: {loading ? '...' : `$${balance.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}
              </div>
            </section>

            <section className={s.cardStat}>
              <div className={s.caption}>Общий PnL (от $10,000)</div>
              <div className={pnlAbs >= 0 ? s.valueGreen : s.valueRed}>
                {loading ? '...' : `${pnlPct >= 0 ? '+' : ''}${pnlPct.toFixed(2)}%`}
              </div>
              <div className={pnlAbs >= 0 ? s.badgeGreen : s.badgeRed}>
                {loading ? '...' : `${pnlAbs >= 0 ? '+' : ''}$${Math.abs(pnlAbs).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}
              </div>
            </section>

            <section className={s.cardBest}>
              <div className={s.cardTitle}>Лучшие монеты</div>
              <div className={s.bestGrid}>
                {bestCoins.length > 0 ? (
                  bestCoins.map(c => (
                    <div key={c.symbol} className={s.bestItem}>
                      <div className={s.bestTitle}>{c.title}</div>
                      <div className={s.bestPnl}>
                        PnL: <span className={c.pnlPct >= 0 ? s.textGreen : s.textRed}>
                          {c.pnlPct >= 0 ? '+' : ''}{c.pnlPct.toFixed(2)}%
                        </span>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className={s.emptyBest}>Нет активных монет</div>
                )}
              </div>
            </section>
          </main>
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  userId: state.user.userId
});

export default connect(mapStateToProps)(Profile);
