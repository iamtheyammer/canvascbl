import thunk from 'redux-thunk';

import { createStore, applyMiddleware, compose } from 'redux';
import createSagaMiddleware from 'redux-saga';
import env from '../util/env';

import canvasRootSaga from '../sagas/canvas';
import reducers from '../reducers';

const sagaMiddleware = createSagaMiddleware();

const middlewares = [thunk, sagaMiddleware];

const composer =
  env.nodeEnv === 'development'
    ? window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose
    : compose;

const store = createStore(reducers, composer(applyMiddleware(...middlewares)));

sagaMiddleware.run(canvasRootSaga);

export default store;
