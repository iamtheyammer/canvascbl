import React from 'react';
import { connect } from 'react-redux';
import v4 from 'uuid/v4';
import { throttle } from 'lodash';
import { isMobile } from 'react-device-detect';
import { Switch, Typography } from 'antd';
import {
  Accordion as MobileAccordion,
  List as MobileList,
  Switch as MobileSwitch
} from 'antd-mobile';
import {
  destinationNames,
  destinationTypes,
  itemTypes,
  tableNames,
  trackCourseVisibilityToggle,
  TrackingLink,
  trackTableRowExpansion,
  vias
} from '../../../util/tracking';
import Padding from '../../Padding';
import { toggleCourseVisibility } from '../../../actions/canvas';
import { setToggleCourseVisibilityId } from '../../../actions/components/coursesettings';
import PopoutLink from '../../PopoutLink';

function CourseSettings(props) {
  const {
    record,
    loading,
    toggleVisibilityIds,
    toggleVisibilityErrors,
    dispatch
  } = props;

  const toggleVisibilityErr =
    toggleVisibilityErrors && toggleVisibilityErrors[record.id];

  function toggleVisibility(toggle) {
    const id = v4();
    dispatch(toggleCourseVisibility(id, record.id, toggle));
    dispatch(setToggleCourseVisibilityId(id, record.id));
    trackCourseVisibilityToggle(record.id, !toggle);
  }

  const switchState = record.canvascblHidden;
  const switchOnChange = throttle(toggle => toggleVisibility(toggle), 2000, {
    leading: true
  });
  const switchDisabled = !!toggleVisibilityErr;
  const switchLoading =
    toggleVisibilityIds &&
    toggleVisibilityIds[record.id] &&
    loading.includes(toggleVisibilityIds[record.id]);

  if (isMobile) {
    return (
      <MobileAccordion
        onChange={open =>
          trackTableRowExpansion(
            tableNames.grades.grades,
            record.id,
            itemTypes.course,
            open.length > 0,
            record.id
          )
        }
      >
        <MobileAccordion.Panel header="Course Settings">
          <MobileList>
            <MobileList.Item
              extra={
                <MobileSwitch
                  checked={switchState}
                  onChange={switchOnChange}
                  disabled={switchDisabled || switchLoading}
                />
              }
              style={{ marginLeft: '10px' }}
              multipleLine
            >
              Hide This Course{toggleVisibilityErr && ' [ERROR!]'}
              <MobileList.Item.Brief>
                Learn more about
                <br />
                hiding courses
                <br />{' '}
                <PopoutLink
                  url="https://go.canvascbl.com/help/hiding-courses"
                  tracking={{
                    destinationName: destinationNames.helpdesk,
                    destinationType: destinationTypes.helpdesk.hidingCourses,
                    via: vias.courseSettingsHideThisCourseLearnMoreLink
                  }}
                  addIcon
                >
                  here
                </PopoutLink>
                .
              </MobileList.Item.Brief>
            </MobileList.Item>
          </MobileList>
        </MobileAccordion.Panel>
      </MobileAccordion>
    );
  }

  return (
    <>
      <Typography.Title level={3}>Course Settings</Typography.Title>
      <Typography.Title level={4}>Hide This Course</Typography.Title>
      <Typography.Text>
        Hiding a course hides it from the Grades table in CanvasCBL. It does not
        hide it in Canvas. If you would like to un-hide a course, go to the{' '}
        <TrackingLink
          to="/dashboard/settings"
          via={vias.courseSettingsHideAClassSettingsLink}
        >
          Settings page
        </TrackingLink>{' '}
        and enable 'Show hidden courses'. Learn more about hiding courses{' '}
        <PopoutLink
          url="https://go.canvascbl.com/help/hiding-courses"
          tracking={{
            destinationName: destinationNames.helpdesk,
            destinationType: destinationTypes.helpdesk.hidingCourses,
            via: vias.courseSettingsHideThisCourseLearnMoreLink
          }}
          addIcon
        >
          here
        </PopoutLink>
        .
      </Typography.Text>
      <Padding all={5} />
      {toggleVisibilityErr && (
        <>
          <Typography.Text type="danger">
            There was an error toggling this course's visibility. Please try
            again later.
          </Typography.Text>
          <Padding all={5} />
        </>
      )}
      <Switch
        checked={switchState}
        onChange={switchOnChange}
        disabled={switchDisabled}
        loading={switchLoading}
      />
    </>
  );
}

export default connect(state => ({
  loading: state.loading,
  courses: state.canvas.courses,
  toggleVisibilityIds: state.components.coursesettings.toggleVisibilityIds,
  toggleVisibilityErrors: state.canvas.toggleCourseVisibilityErrors
}))(CourseSettings);
