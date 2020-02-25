import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { Typography, Switch, Skeleton } from 'antd';
import { throttle } from 'lodash';
import v4 from 'uuid';
import Padding from '../../Padding';
import {
  getNotificationSettingsAndTypes,
  toggleNotificationType
} from '../../../actions/settings';
import {
  setGetNotificationSettingsAndTypesId,
  setToggleNotificationStatusId
} from '../../../actions/components/settings';

function Settings(props) {
  const {
    settingsAndTypesId,
    settingsAndTypesError,
    toggleIds,
    notifications,
    loading,
    dispatch
  } = props;

  function fetchSettingsAndTypes() {
    const id = v4();
    dispatch(getNotificationSettingsAndTypes(id));
    dispatch(setGetNotificationSettingsAndTypesId(id));
  }

  if (
    !settingsAndTypesId &&
    (!notifications || !notifications.settings || !notifications.types)
  ) {
    fetchSettingsAndTypes();
  }

  function toggleNotification(typeId, toggle) {
    const id = v4();
    dispatch(toggleNotificationType(id, typeId, toggle));
    dispatch(setToggleNotificationStatusId(id, typeId));
  }

  const notificationSettingsStatus = {
    loading:
      // prevents skeleton from showing on refetches
      (!toggleIds && loading.includes(settingsAndTypesId)) ||
      !notifications ||
      !notifications.settings ||
      !notifications.types,
    disabled: !!settingsAndTypesError
  };

  return (
    <>
      <Typography.Title level={2}>Settings</Typography.Title>
      <Typography.Text type="secondary">
        Configure CanvasCBL to your liking.
      </Typography.Text>
      <Typography.Title level={3}>Notifications</Typography.Title>
      <Typography.Text type="secondary">
        Currently, all notifications are delivered via email.
      </Typography.Text>
      {!!settingsAndTypesError ? (
        <>
          <Padding all={10} />
          <Typography.Text type="danger">
            There was an error fetching notification settings. Please contact
            support or try again later.
          </Typography.Text>
        </>
      ) : !notificationSettingsStatus.loading ? (
        notifications.types.map(t => (
          <Fragment key={t.short_name}>
            <Typography.Title level={4}>{t.name}</Typography.Title>
            <Typography.Text>{t.description}</Typography.Text>
            <Padding all={5} />
            {notifications.toggleErrors && notifications.toggleErrors[t.id] && (
              <>
                <Typography.Text type="danger">
                  There was an error toggling this notification. Please contact
                  support or try again later.
                </Typography.Text>
                <Padding all={5} />
              </>
            )}
            <Switch
              onChange={throttle(
                checked => toggleNotification(t.id, checked),
                2000
              )}
              checked={
                notifications &&
                !!notifications.settings.filter(
                  ns => ns.notification_type_id === t.id
                )[0]
              }
              loading={
                notificationSettingsStatus.loading ||
                (toggleIds &&
                  toggleIds[t.id] &&
                  loading.includes(toggleIds[t.id]))
              }
              disabled={
                notificationSettingsStatus.disabled ||
                (notifications.toggleErrors && notifications.toggleErrors[t.id])
              }
            />
          </Fragment>
        ))
      ) : (
        <Skeleton active />
      )}
    </>
  );
}

export default connect(state => ({
  settingsAndTypesId:
    state.components.settings.getNotificationSettingsAndTypesId,
  settingsAndTypesError: state.settings.getNotificationSettingsAndTypesError,
  toggleIds: state.components.settings.toggleIds,
  notifications: state.settings.notifications,
  loading: state.loading
}))(Settings);
