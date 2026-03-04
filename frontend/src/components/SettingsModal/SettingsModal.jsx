import React, { useState } from 'react';
import { connect } from 'react-redux';
import s from './SettingsModal.module.css';
import { authApi } from '../../api/auth';
import { resetPortfolio } from '../../api/portfolio';
import { resetOrders } from '../../api/trading';
import { resetBets } from '../../api/bets';

const SettingsModal = ({ isOpen, onClose, userId }) => {
  const [loginForm, setLoginForm] = useState({ newLogin: '', message: '', error: '' });
  const [passForm, setPassForm] = useState({ oldPass: '', newPass: '', message: '', error: '' });
  const [showResetConfirm, setShowResetConfirm] = useState(false);
  const [resetStatus, setResetStatus] = useState({ message: '', error: '' });

  if (!isOpen) return null;

  const handleLoginChange = async (e) => {
    e.preventDefault();
    setLoginForm({ ...loginForm, message: '', error: '' });
    try {
      await authApi.changeLogin(loginForm.newLogin);
      setLoginForm({ newLogin: '', message: 'Логин успешно обновлен', error: '' });
    } catch (err) {
      setLoginForm({ ...loginForm, error: err.response?.data?.error || 'Не удалось обновить логин' });
    }
  };

  const handlePassChange = async (e) => {
    e.preventDefault();
    setPassForm({ ...passForm, message: '', error: '' });
    try {
      await authApi.changePassword(passForm.oldPass, passForm.newPass);
      setPassForm({ oldPass: '', newPass: '', message: 'Пароль успешно обновлен', error: '' });
    } catch (err) {
      setPassForm({ ...passForm, error: err.response?.data?.error || 'Не удалось обновить пароль' });
    }
  };

  const handleReset = async () => {
    setResetStatus({ message: 'Сброс...', error: '' });
    try {
      await resetPortfolio(userId);
      await resetOrders(userId);
      await resetBets(userId);

      setResetStatus({ message: 'Аккаунт успешно сброшен. Баланс восстановлен до 10,000 USDT.', error: '' });
      setShowResetConfirm(false);

      setTimeout(() => {
        window.location.reload();
      }, 2000);
    } catch (err) {
      console.error(err);
      setResetStatus({ message: '', error: 'Не удалось сбросить аккаунт. Попробуйте снова.' });
    }
  };

  return (
    <div className={s.modalOverlay} onClick={onClose}>
      <div className={s.modalContent} onClick={e => e.stopPropagation()}>
        <button className={s.closeBtn} onClick={onClose}>&times;</button>
        <h2 className={s.title}>Настройки</h2>

        <div className={s.section}>
          <h3 className={s.sectionTitle}>Сменить логин</h3>
          <form onSubmit={handleLoginChange}>
            <div className={s.formGroup}>
              <label className={s.label}>Новый логин</label>
              <input
                type="text"
                className={s.input}
                value={loginForm.newLogin}
                onChange={e => setLoginForm({ ...loginForm, newLogin: e.target.value })}
                required
              />
            </div>
            <button type="submit" className={`${s.btn} ${s.btnPrimary}`}>Обновить логин</button>
            {loginForm.message && <div className={`${s.message} ${s.success}`}>{loginForm.message}</div>}
            {loginForm.error && <div className={`${s.message} ${s.error}`}>{loginForm.error}</div>}
          </form>
        </div>

        <div className={s.section}>
          <h3 className={s.sectionTitle}>Сменить пароль</h3>
          <form onSubmit={handlePassChange}>
            <div className={s.formGroup}>
              <label className={s.label}>Старый пароль</label>
              <input
                type="password"
                className={s.input}
                value={passForm.oldPass}
                onChange={e => setPassForm({ ...passForm, oldPass: e.target.value })}
                required
              />
            </div>
            <div className={s.formGroup}>
              <label className={s.label}>Новый пароль</label>
              <input
                type="password"
                className={s.input}
                value={passForm.newPass}
                onChange={e => setPassForm({ ...passForm, newPass: e.target.value })}
                required
              />
            </div>
            <button type="submit" className={`${s.btn} ${s.btnPrimary}`}>Обновить пароль</button>
            {passForm.message && <div className={`${s.message} ${s.success}`}>{passForm.message}</div>}
            {passForm.error && <div className={`${s.message} ${s.error}`}>{passForm.error}</div>}
          </form>
        </div>

        <div className={s.section}>
          <h3 className={s.sectionTitle}>Опасная зона</h3>
          {!showResetConfirm ? (
            <button
              className={`${s.btn} ${s.btnDanger}`}
              onClick={() => setShowResetConfirm(true)}
            >
              Сбросить аккаунт
            </button>
          ) : (
            <div className={s.resetConfirm}>
              <p className={s.resetText}>
                Вы уверены? Это действие:
                <br />• Сбросит ваш баланс до 10,000 USDT
                <br />• Удалит все ваши ордера и ставки
                <br />• Очистит историю портфеля
                <br /><b>Это действие нельзя отменить.</b>
              </p>
              <div className={s.confirmActions}>
                <button
                  className={`${s.btn} ${s.btnDanger}`}
                  onClick={handleReset}
                >
                  Подтвердить сброс
                </button>
                <button
                  className={`${s.btn} ${s.btnCancel}`}
                  onClick={() => setShowResetConfirm(false)}
                >
                  Отмена
                </button>
              </div>
            </div>
          )}
          {resetStatus.message && <div className={`${s.message} ${s.success}`}>{resetStatus.message}</div>}
          {resetStatus.error && <div className={`${s.message} ${s.error}`}>{resetStatus.error}</div>}
        </div>

      </div>
    </div>
  );
};

const mapStateToProps = (state) => ({
  userId: state.user.userId
});

export default connect(mapStateToProps)(SettingsModal);
