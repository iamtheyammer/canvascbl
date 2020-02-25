import { all, put, takeLeading } from 'redux-saga/effects';
import { endLoading, startLoading } from '../actions/loading';
import makeApiRequest from '../util/api/makeApiRequest';
import {
  getNotificationSettingsAndTypesError,
  gotNotificationSettingsAndTypes,
  SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES,
  SETTINGS_TOGGLE_NOTIFICATION_TYPE,
  toggledNotificationType,
  toggleNotificationTypeError
} from '../actions/settings';

function* fetchNotificationSettingsAndTypes({ id }) {
  yield put(startLoading(id));
  try {
    const settingsReq = yield makeApiRequest('notifications/settings', {
      include: ['notification_types']
    });
    yield put(
      gotNotificationSettingsAndTypes(
        settingsReq.data.notification_settings,
        settingsReq.data.notification_types
      )
    );
  } catch (e) {
    yield put(
      getNotificationSettingsAndTypesError(
        e.response ? e.response.data : "can't connect"
      )
    );
  }
  yield put(endLoading(id));
}

function* toggleNotificationType({ id, toggle, typeId }) {
  yield put(startLoading(id));
  try {
    yield makeApiRequest(
      `notifications/types/${typeId}`,
      {
        medium: 'email'
      },
      toggle ? 'put' : 'delete'
    );
    yield put(toggledNotificationType(typeId));
    yield fetchNotificationSettingsAndTypes(id);
  } catch (e) {
    yield put(
      toggleNotificationTypeError(
        e.response ? e.response.data : "can't connect",
        typeId
      )
    );
  }
  yield put(endLoading(id));
}

function* watcher() {
  yield takeLeading(
    SETTINGS_GET_NOTIFICATION_SETTINGS_AND_TYPES,
    fetchNotificationSettingsAndTypes
  );
  yield takeLeading(SETTINGS_TOGGLE_NOTIFICATION_TYPE, toggleNotificationType);
}

export default function* settingsRootSaga() {
  yield all([watcher()]);
}
