import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import ConnectedApp from './App';
import { CookiesProvider } from 'react-cookie';
import store from './store/index';

// import * as serviceWorker from './serviceWorker';

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
