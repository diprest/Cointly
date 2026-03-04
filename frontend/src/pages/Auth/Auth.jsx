import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { toast } from 'react-toastify';
import s from './Auth.module.css';
import Header from '../../components/Header/Header';
import { authApi } from '../../api/auth';
import { login } from '../../store/userSlice';

export default function Auth() {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();

  const [isRegister, setIsRegister] = useState(false);
  const [loginVal, setLoginVal] = useState('');
  const [pass, setPass] = useState('');
  const [confirmPass, setConfirmPass] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const from = location.state?.from?.pathname || '/';

  const handleSubmit = async () => {
    if (!loginVal || !pass) {
      toast.error('Заполните все поля');
      return;
    }

    if (isRegister && pass !== confirmPass) {
      toast.error('Пароли не совпадают');
      return;
    }

    setIsLoading(true);
    try {
      let data;
      if (isRegister) {
        data = await authApi.register(loginVal, pass);
        toast.success('Регистрация успешна! Теперь войдите.');
        setIsRegister(false);
        setIsLoading(false);
        return;
      } else {
        data = await authApi.login(loginVal, pass);
      }

      if (data && data.token) {
        dispatch(login({ token: data.token, userId: data.user_id }));
        toast.success('Вход выполнен!');
        navigate(from, { replace: true });
      } else {
        toast.error('Ошибка: не получен токен');
      }
    } catch (err) {
      console.error(err);
      toast.error(err.response?.data?.error || 'Ошибка авторизации');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={s.wrap}>
      <Header />

      <main className={s.main}>
        <section className={s.card}>
          <h2 className={s.title}>
            {isRegister ? 'Регистрация' : 'Вход'}
          </h2>

          <label className={s.field}>
            <input
              className={s.input}
              type="text"
              placeholder="Логин"
              value={loginVal}
              onChange={(e) => setLoginVal(e.target.value)}
            />
          </label>

          <label className={s.field}>
            <input
              className={s.input}
              type="password"
              placeholder="Пароль"
              value={pass}
              onChange={(e) => setPass(e.target.value)}
            />
          </label>

          {isRegister && (
            <label className={s.field}>
              <input
                className={s.input}
                type="password"
                placeholder="Подтвердите пароль"
                value={confirmPass}
                onChange={(e) => setConfirmPass(e.target.value)}
              />
            </label>
          )}

          <div className={s.actions}>
            <button
              className={s.btnPrimary}
              onClick={handleSubmit}
              disabled={isLoading}
            >
              {isLoading ? 'Загрузка...' : (isRegister ? 'Создать аккаунт' : 'Войти')}
            </button>

            <div className={s.switchWrap}>
              <span className={s.switchText}>
                {isRegister ? 'Уже есть аккаунт?' : 'Ещё нет аккаунта?'}
              </span>
              <button
                className={s.btnLink}
                onClick={() => setIsRegister(!isRegister)}
              >
                {isRegister ? 'Войти' : 'Создать аккаунт'}
              </button>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}