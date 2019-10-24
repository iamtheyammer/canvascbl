import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import ConnectedApp from './App';
import { createStore } from 'redux';
import { CookiesProvider } from 'react-cookie';
import reducers from './reducers';
import middlewares from './middleware/index';

// import * as serviceWorker from './serviceWorker';

const store = createStore(reducers, middlewares);

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
