import { createStore, applyMiddleware, compose } from 'redux';
import createSagaMiddleware from 'redux-saga';
import env from '../util/env';

import reducers from '../reducers';
import rootSaga from '../sagas';

const sagaMiddleware = createSagaMiddleware();

const middlewares = [sagaMiddleware];

const composer =
  env.nodeEnv === 'development'
    ? window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose
    : compose;

const store = createStore(reducers, composer(applyMiddleware(...middlewares)));

sagaMiddleware.run(rootSaga);

export default store;
