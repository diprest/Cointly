import React, { useEffect, useState, useRef } from 'react';
import { useSelector } from 'react-redux';
import { getBalance } from '../../api/portfolio';
import { createOrder, getOrders, cancelOrder } from '../../api/trading';
import { createBet, getUserBets } from '../../api/bets';
import { getTicker } from '../../api/market';
import { useParams, Link } from 'react-router-dom';
import { toast } from 'react-toastify';
import s from './Coin.module.css';
import TradingViewWidget from '../../components/TradingViewWidget';
import OrderForm from './components/OrderForm';
import BettingWidget from './components/BettingWidget';
import ActiveOrders from './components/ActiveOrders';
import ActiveBets from './components/ActiveBets';

import Header from '../../components/Header/Header';

function useSymbol() {
  const { symbol } = useParams();
  let s = (symbol || 'BTCUSDT').toUpperCase();
  if (!s.endsWith('USDT')) {
    s += 'USDT';
  }
  return s;
}

export default function Coin() {
  const symbol = useSymbol();
  const [series, setSeries] = useState([]);
  const { userId } = useSelector(state => state.user);
  const [balance, setBalance] = useState(0);
  const [coinBalance, setCoinBalance] = useState(0);

  const [activeOrders, setActiveOrders] = useState([]);
  const [activeBets, setActiveBets] = useState([]);

  const [lastPrice, setLastPrice] = useState(0);

  const prevOrdersRef = useRef([]);
  const prevBetsRef = useRef([]);

  const checkForOrderUpdates = (prev, current) => {
    if (!prev || prev.length === 0) return;

    current.forEach(curr => {
      const old = prev.find(p => p.id === curr.id);
      if (old && old.status === 'NEW' && curr.status === 'FILLED') {
        const action = curr.side === 'BUY' ? 'Покупка' : 'Продажа';
        const base = curr.symbol.replace('USDT', '');
        toast.success(`${action} исполнена! ${curr.amount} ${base} по цене ${curr.price}`);
      }
    });
  };

  const fetchData = async () => {
    try {
      const ticker = await getTicker(symbol);
      setLastPrice(ticker.price);
    } catch (e) {
      console.error("Ticker error", e);
    }

    if (!userId) return;

    try {
      const balData = await getBalance(userId, 'USDT');
      setBalance(Number(balData.amount) - Number(balData.locked_bal || 0));
    } catch (e) {
      console.error("Balance USDT error", e);
    }

    try {
      const baseAsset = symbol.replace('USDT', '');
      const coinBalData = await getBalance(userId, baseAsset);
      setCoinBalance(Number(coinBalData.amount) - Number(coinBalData.locked_bal || 0));
    } catch (e) {
      console.error("Balance Coin error", e);
    }

    try {
      const orders = await getOrders(userId);
      checkForOrderUpdates(prevOrdersRef.current, orders);
      prevOrdersRef.current = orders;
      setActiveOrders(orders.filter(o => o.symbol === symbol && o.status === 'NEW'));
    } catch (e) {
      console.error("Orders error", e);
    }

    try {
      const bets = await getUserBets(userId);
      prevBetsRef.current = bets;
      setActiveBets(bets.filter(b => b.symbol === symbol && b.status === 'OPEN'));
    } catch (e) {
      console.error("Bets error", e);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, [userId, symbol]);

  const handleOrder = async (side, qty, price) => {
    if (!userId) {
      toast.error("Пожалуйста, войдите в систему");
      return;
    }
    const cost = Number(qty) * Number(price || lastPrice);
    if (side.toUpperCase() === 'BUY' && cost > balance) {
      toast.error(`Недостаточно средств. Требуется: ${cost.toFixed(2)} USDT, Доступно: ${balance.toFixed(2)} USDT`);
      return;
    }
    if (side.toUpperCase() === 'SELL' && Number(qty) > coinBalance) {
      toast.error(`Недостаточно монет. Требуется: ${qty} ${symbol.replace('USDT', '')}, Доступно: ${coinBalance}`);
      return;
    }

    try {
      const order = await createOrder({
        user_id: Number(userId),
        symbol,
        side: side.toUpperCase(),
        type: price ? 'LIMIT' : 'MARKET',
        amount: Number(qty),
        price: Number(price || lastPrice)
      });

      if (order.status === 'FILLED') {
        const action = side.toUpperCase() === 'BUY' ? 'Покупка' : 'Продажа';
        const base = symbol.replace('USDT', '');
        toast.success(`${action} исполнена! ${order.amount} ${base} по цене ${order.price}`);
      } else {
        toast.success(`Ордер на ${side === 'buy' ? 'покупку' : 'продажу'} создан!`);
      }

      fetchData();
    } catch (e) {
      console.error(e);
      toast.error("Ошибка создания ордера: " + (e.response?.data?.error || e.message));
    }
  };


  const handleBet = async (direction, amount, duration) => {
    if (!userId) {
      toast.error("Пожалуйста, войдите в систему");
      return;
    }
    try {
      await createBet({
        user_id: Number(userId),
        symbol: symbol,
        direction: direction === 'up' ? 'UP' : 'DOWN',
        stake_amount: Number(amount),
        duration_sec: Number(duration)
      });
      toast.success(`Ставка ${direction === 'up' ? 'ВВЕРХ' : 'ВНИЗ'} принята!`);
      fetchData();
    } catch (e) {
      console.error(e);
      toast.error("Ошибка ставки: " + (e.response?.data?.error || e.message));
    }
  };

  const handleCancelOrder = async (id) => {
    try {
      await cancelOrder(id, userId);
      toast.info("Ордер отменен");
      fetchData();
    } catch (e) {
      console.error(e);
      toast.error("Не удалось отменить ордер");
    }
  };

  return (
    <div className={s.wrap}>
      <Header>
        <Link to="/profile" className={s.profileBtn}>← Профиль</Link>
      </Header>

      <div className={s.main}>
        <div className={s.grid}>
          <div className={s.contentCol}>
            <div className={s.chartSection}>
              <div className={s.chartInner}>
                <TradingViewWidget symbol={symbol} />
              </div>
            </div>

            <div className={s.listsGrid}>
              <ActiveOrders orders={activeOrders} onCancel={handleCancelOrder} />
              <ActiveBets bets={activeBets} />
            </div>
          </div>

          <div className={s.sidebarCol}>
            <div className={s.balanceBlock}>
              <div className={s.balanceRow}>
                <span className={s.balanceLabel}>Доступно:</span>
                <span className={s.balanceValue}>{Number(balance).toFixed(2)} USDT</span>
              </div>
              <div className={s.balanceRow}>
                <span className={s.balanceLabel}>В наличии:</span>
                <span className={s.balanceValue}>{Number(coinBalance).toFixed(6)} {symbol.replace('USDT', '')}</span>
              </div>
            </div>

            <OrderForm lastPrice={lastPrice} onOrder={handleOrder} />
            <BettingWidget onBet={handleBet} />
          </div>
        </div>
      </div>
    </div>
  );
}
