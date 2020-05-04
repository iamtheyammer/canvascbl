import React, { Fragment, useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Typography, Switch, Skeleton } from 'antd';
import { throttle } from 'lodash';
import v4 from 'uuid';
import Padding from '../../Padding';
import {
  getNotificationSettingsAndTypes,
  toggledShowHiddenCourses,
  toggleNotificationType
} from '../../../actions/settings';
import {
  setGetNotificationSettingsAndTypesId,
  setToggleNotificationStatusId
} from '../../../actions/components/settings';
import {
  destinationNames,
  destinationTypes,
  pageNames,
  trackChangedHiddenCourseVisibility,
  trackNotificationStatusToggle,
  trackPageView,
  vias
} from '../../../util/tracking';
import PopoutLink from '../../PopoutLink';

function Settings(props) {
  const {
    settingsAndTypesId,
    settingsAndTypesError,
    toggleIds,
    showHiddenCourses,
    notifications,
    loading,
    dispatch
  } = props;

  const [loaded, setLoaded] = useState(false);
  useEffect(() => {
    /*
    This system is to prevent sending tons of Page View events to Mixpanel.
    Those tons of events are sent because, every time the state changes,
    this component is rerendered. The most common state change is when
    a grade average loads in for plus users.

    It works with two hooks: state and effect.

    There's a loaded state hook set to false just above.

    The effect hook is used here to run whenever loaded changes--
    if it's true, we'll track a page view. If not, whatever.

    The reason that this works is because state is reset on unmount.
    So we only get one page view per actual page view.
     */

    if (loaded) {
      trackPageView(pageNames.settings);
    }
  }, [loaded]);

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

  function toggleNotification(typeId, typeShortName, toggle) {
    const id = v4();
    dispatch(toggleNotificationType(id, typeId, toggle));
    dispatch(setToggleNotificationStatusId(id, typeId));

    trackNotificationStatusToggle(toggle, typeShortName);
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

  if (!loaded) {
    setLoaded(true);
  }
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
        notifications.types.map((t) => (
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
                (checked) => toggleNotification(t.id, t.short_name, checked),
                2000
              )}
              checked={
                notifications &&
                !!notifications.settings.filter(
                  (ns) => ns.notification_type_id === t.id
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
      <Padding all={10} />
      <Typography.Title level={3}>Hidden Courses</Typography.Title>
      <Typography.Text>
        You can hide a course by expanding a course to view its settings, then
        enabling Hide This Course. Learn more about hiding courses{' '}
        <PopoutLink
          url="https://go.canvascbl.com/help/hiding-courses"
          tracking={{
            destinationName: destinationNames.helpdesk,
            destinationType: destinationTypes.helpdesk.hidingCourses,
            via: vias.settingsShowHiddenCoursesDescriptionLearnMoreLink
          }}
          addIcon
        >
          here
        </PopoutLink>
        . Do you want to show hidden courses?
      </Typography.Text>
      <Padding all={5} />
      <Switch
        onChange={(toggle) => {
          dispatch(toggledShowHiddenCourses(toggle));
          trackChangedHiddenCourseVisibility(toggle);
        }}
        checked={showHiddenCourses}
      />
    </>
  );
}

export default connect((state) => ({
  settingsAndTypesId:
    state.components.settings.getNotificationSettingsAndTypesId,
  settingsAndTypesError: state.settings.getNotificationSettingsAndTypesError,
  toggleIds: state.components.settings.toggleIds,
  showHiddenCourses: state.settings.showHiddenCourses,
  notifications: state.settings.notifications,
  loading: state.loading
}))(Settings);
