import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import ConnectedApp from './App';
import { CookiesProvider } from 'react-cookie';
import store from './store/index';

// import * as serviceWorker from './serviceWorker';

console.log(
  '---------------\n',
  'Hey there, curious fellow! While we love inquisitive minds, please note that reverse-engineering any part of CanvasCBL is against our terms of service. Thank you for using CanvasCBL!\n',
  '---------------'
);

ReactDOM.render(
  <CookiesProvider>
    <Provider store={store}>
      <ConnectedApp />
    </Provider>
  </CookiesProvider>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
// serviceWorker.unregister();
