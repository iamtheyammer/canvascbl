import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import v4 from 'uuid/v4';

import {
  Typography,
  Table,
  Icon,
  Tag,
  Skeleton,
  Popover,
  Row,
  Col,
  Statistic,
  Divider
} from 'antd';
import { Accordion as MobileAccordion, List as MobileList } from 'antd-mobile';

import { gradeMapByGrade } from '../../../util/canvas/gradeMapByGrade';
import getActiveCourses from '../../../util/canvas/getActiveCourses';

import env from '../../../util/env';
import { ReactComponent as PopOutIcon } from '../../../assets/pop_out.svg';
import { ReactComponent as plusIcon } from '../../../assets/plus.svg';
import { desc } from '../../../util/sort';
import PopoutLink from '../../PopoutLink';
import { getPreviousGrades } from '../../../actions/plus';
import moment from 'moment';
import { isMobile } from 'react-device-detect';
import Loading from '../Loading';
import Padding from '../../Padding';
import ConnectedObserveeHandler from '../DashboardNav/ObserveeHandler';
import roundNumberToDigits from '../../../util/roundNumberToDigits';
import {
  destinationNames,
  destinationTypes,
  itemTypes,
  pageNames,
  tableNames,
  TrackingLink,
  trackPageView,
  trackTableRowExpansion,
  vias
} from '../../../util/tracking';
import { truncate } from 'lodash';
import CourseSettings from './CourseSettings';

function PreviousGrade(props) {
  const { userHasValidSubscription, grade, previousGrade } = props;

  if (!userHasValidSubscription) {
    return (
      <Popover
        title="CanvasCBL+ Required"
        content="CanvasCBL+ is required to use this feature. Go to the Upgrades page to upgrade!"
      >
        <Icon component={plusIcon} /> Required
      </Popover>
    );
  }

  if (previousGrade === undefined) return <Tag>Unavailable</Tag>;

  if (previousGrade === 'loading') {
    return <Skeleton paragraph={false} active title={{ width: '50%' }} />;
  }

  const prevGrade = gradeMapByGrade[previousGrade.grade];
  const currentGrade = gradeMapByGrade[grade];
  if (!prevGrade || !currentGrade) return <Tag>Unavailable</Tag>;

  let color = '';

  // old is better than new
  if (prevGrade.rank > currentGrade.rank) {
    color = 'volcano';
  } else if (prevGrade.rank < currentGrade.rank) {
    color = 'green';
  }

  return (
    <Popover
      title={`Previous Grade: ${previousGrade.grade}`}
      content={`From: ${moment.unix(previousGrade.insertedAt).calendar()}`}
    >
      <Tag color={color}>{previousGrade.grade}</Tag>
    </Popover>
  );
}

const tableColumns = [
  {
    title: 'Class Name',
    dataIndex: 'name',
    key: 'name',
    sorter: (a, b) => desc(a.name, b.name),
    render: (text, record) =>
      record.breakdownIsAvailable ? (
        <TrackingLink
          to={`/dashboard/grades/${record.id}`}
          pageName={pageNames.gradeBreakdown}
          via={vias.gradesTableCourseName}
        >
          {text}
          {record.canvascblHidden && (
            <Popover
              title="Hidden Course"
              content={
                'This course is normally hidden, but you have Show Hidden Courses enabled.'
              }
              placement="topLeft"
            >
              <Divider type="vertical" />
              <Icon type="eye-invisible" />
            </Popover>
          )}
        </TrackingLink>
      ) : (
        <>
          {text}

          {record.canvascblHidden && (
            <Popover
              title="Hidden Course"
              content={
                'This course is normally hidden, but you have Show Hidden Courses enabled.'
              }
              placement="topLeft"
            >
              <Divider type="vertical" />
              <Icon type="eye-invisible" />
            </Popover>
          )}
        </>
      )
  },
  {
    title: 'Grade',
    dataIndex: 'grade',
    key: 'grade',
    sorter: (a, b) => desc(a.grade, b.grade),
    defaultSortOrder: 'desc',
    render: (text) => <Typography.Text strong>{text}</Typography.Text>
  },
  {
    title: (
      <Popover
        title="Grades From Last Login"
        content="Hover over a previous grade to see when it's from."
      >
        <Icon component={plusIcon} /> Previous Grade
      </Popover>
    ),
    dataIndex: 'averageGrade',
    key: 'averageGrade',
    render: (text, record) => (
      <PreviousGrade
        userHasValidSubscription={record.userHasValidSubscription}
        grade={record.grade}
        previousGrade={record.previousGrade}
      />
    )
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (text, record) => (
      <>
        {record.breakdownIsAvailable && (
          <>
            <TrackingLink
              to={`/dashboard/grades/${record.id}`}
              pageName={pageNames.gradeBreakdown}
              via={vias.gradesTableSeeBreakdownLink}
            >
              See Breakdown
            </TrackingLink>
            <Divider type="vertical" />
          </>
        )}
        <PopoutLink
          url={`https://${env.defaultSubdomain}.instructure.com/courses/${record.id}`}
          tracking={{
            destinationName: destinationNames.canvas,
            destinationType: destinationTypes.course,
            via: vias.gradesTableOpenOnCanvas
          }}
          addIcon
        >
          Open on Canvas
        </PopoutLink>
      </>
    )
  }
];

function Grades(props) {
  const [getPrevGradeId, setGetPrevGradeId] = useState('');

  const {
    dispatch,
    grades,
    subdomain,
    showHiddenCourses,
    loading,
    error,
    courses,
    users,
    allGpas,
    activeUserId,
    observees,
    plus
  } = props;

  const err = error[getPrevGradeId];

  const activeUser = users && activeUserId && users[activeUserId];

  useEffect(() => {
    if (
      activeUser &&
      plus.session &&
      plus.session.has_valid_subscription &&
      !plus.previousGrades &&
      !getPrevGradeId
    ) {
      const id = v4();
      dispatch(getPreviousGrades(id));
      setGetPrevGradeId(id);
    }
    // ignoring because we only want this hook to re-run on a prop change
    // eslint-disable-next-line
  }, [props]);

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
      trackPageView(pageNames.grades);
    }
  }, [loaded]);

  if (err) {
    return (
      <Typography.Text type={'danger'}>
        There was an unknown error getting your grades. Please try again later.
      </Typography.Text>
    );
  }
  if (!activeUser || !plus.session || !courses) {
    return <Loading text="grades" />;
  }

  const activeCourses =
    courses && activeUser ? getActiveCourses(courses, activeUser.id) : [];

  const previousGrades =
    plus &&
    plus.previousGrades &&
    plus.previousGrades.filter((pg) => pg.canvasUserId === activeUserId);

  const data = activeCourses.map((c) => {
    const detailedGrade = grades[activeUserId][c.id]
      ? grades[activeUserId][c.id]
      : 'Error, try reloading';

    return {
      key: c.id,
      name: c.name,
      grade: detailedGrade.grade.grade,
      id: c.id,
      canvascblHidden: c.canvascbl_hidden,
      // don't hide if show hidden courses is enabled-- otherwise respect
      hide: !!showHiddenCourses ? false : c.canvascbl_hidden,
      userHasValidSubscription: plus.session.has_valid_subscription,
      breakdownIsAvailable:
        detailedGrade.grade.grade !== 'N/A' &&
        !detailedGrade.grade.grade.toLowerCase().includes('error') &&
        !!detailedGrade.averages,
      previousGrade: loading.includes(getPrevGradeId)
        ? 'loading'
        : plus &&
          plus.previousGrades &&
          !error[getPrevGradeId] &&
          previousGrades.filter((pg) => pg.courseId === c.id)[0] &&
          previousGrades.filter((pg) => pg.courseId === c.id)[0]
    };
  });

  const gradesTitle = (
    <Typography.Title level={2}>
      {observees && observees.length
        ? `${users[activeUserId].name.split(' ')[0]}'s Grades`
        : 'Grades'}
    </Typography.Title>
  );

  const showGpa = !!(
    allGpas &&
    allGpas[activeUserId] &&
    allGpas[activeUserId].unweighted.default !== 0
  );
  const gpa = allGpas && allGpas[activeUserId];

  if (!loaded) {
    setLoaded(true);
  }

  const showData = data.filter((d) => !d.hide);

  if (isMobile) {
    return (
      <>
        {gradesTitle}
        <Padding all={5} />
        <>
          <Typography.Title level={3}>
            Individual Course Grades
          </Typography.Title>
          <MobileAccordion>
            {showData.map((d) => (
              <MobileAccordion.Panel
                key={d.key}
                style={{ padding: '5px 5px 5px 0px' }}
                header={
                  <div style={{ paddingRight: '6px' }}>
                    <div style={{ float: 'left', overflow: 'hidden' }}>
                      {<>{truncate(d.name, { length: 25 })}</>}
                    </div>
                    <div style={{ float: 'right' }}>{d.grade}</div>
                  </div>
                }
              >
                <MobileList style={{ paddingLeft: '6px' }}>
                  {d.canvascblHidden && (
                    <MobileList.Item multipleLine={true}>
                      Hidden Course
                      <MobileList.Item.Brief>
                        This course is normally hidden,
                        <br /> but you have show hidden <br />
                        courses enabled.
                      </MobileList.Item.Brief>
                    </MobileList.Item>
                  )}
                  {d.breakdownIsAvailable && (
                    <MobileList.Item>
                      <TrackingLink
                        to={`/dashboard/grades/${d.id}`}
                        pageName={pageNames.gradeBreakdown}
                        via={vias.gradesTableSeeBreakdownLink}
                      >
                        See Breakdown
                      </TrackingLink>
                    </MobileList.Item>
                  )}
                  <MobileList.Item>
                    <PopoutLink
                      url={`https://${
                        subdomain || 'canvas'
                      }.instructure.com/courses/${d.id}`}
                      tracking={{
                        destinationName: destinationNames.canvas,
                        destinationType: destinationTypes.course,
                        via: vias.gradesTableOpenOnCanvas
                      }}
                    >
                      Open on Canvas <Icon component={PopOutIcon} />
                    </PopoutLink>
                  </MobileList.Item>
                  <MobileList.Item
                    extra={
                      <PreviousGrade
                        userHasValidSubscription={d.userHasValidSubscription}
                        grade={d.grade}
                        previousGrade={d.previousGrade}
                      />
                    }
                  >
                    Previous Grade
                  </MobileList.Item>
                  <CourseSettings record={d} />
                </MobileList>
              </MobileAccordion.Panel>
            ))}
          </MobileAccordion>
        </>
        <Padding all={5} />
        {showGpa && (
          <>
            <Padding all={5} />
            <Typography.Title level={3}>GPA</Typography.Title>
            Unweighted GPAs for the current semester. Learn more{' '}
            <PopoutLink
              url="https://go.canvascbl.com/help/gpas"
              tracking={{
                destinationName: destinationNames.helpdesk,
                destinationType: destinationTypes.helpdesk.gpas,
                via: vias.gpaLearnMore
              }}
            >
              here
            </PopoutLink>
            .
            <Padding all={5} />
            <MobileList>
              <MobileList.Item
                extra={roundNumberToDigits(gpa.unweighted.default, 2)}
              >
                Report Card GPA{' '}
                <PopoutLink
                  url="https://go.canvascbl.com/help/gpas"
                  tracking={{
                    destinationName: destinationNames.helpdesk,
                    destinationType: destinationTypes.helpdesk.gpas,
                    via: vias.gpaReportCardQuestionIcon
                  }}
                >
                  <Icon type="question-circle" />
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item
                extra={roundNumberToDigits(gpa.unweighted.subgrades, 2)}
              >
                Traditional GPA{' '}
                <PopoutLink
                  url="https://go.canvascbl.com/help/gpas"
                  tracking={{
                    destinationName: destinationNames.helpdesk,
                    destinationType: destinationTypes.helpdesk.gpas,
                    via: vias.gpaTraditionalQuestionIcon
                  }}
                >
                  <Icon type="question-circle" />
                </PopoutLink>
              </MobileList.Item>
            </MobileList>
          </>
        )}
        {observees && observees.length > 0 && (
          <>
            <Padding br />
            <Typography.Title level={3}>Switch Students</Typography.Title>
            <ConnectedObserveeHandler via={vias.mobileGradesObserveeSwitcher} />
          </>
        )}
      </>
    );
  }

  return (
    <>
      {gradesTitle}
      <Typography.Text type="secondary">
        If {observees && observees.length ? 'your student has' : 'you have'} a
        grade in a class, click on the name to see a detailed breakdown.
      </Typography.Text>
      <Padding bottom="12px" />
      <>
        <Table
          columns={tableColumns}
          dataSource={showData}
          expandedRowRender={(record) => <CourseSettings record={record} />}
          onExpand={(expanded, record) => {
            trackTableRowExpansion(
              tableNames.grades.grades,
              record.id,
              itemTypes.course,
              expanded,
              record.id
            );
          }}
        />
      </>
      {showGpa && (
        <>
          <Typography.Title level={3}>GPA</Typography.Title>
          <Padding all={5} />
          <Row gutter={16}>
            <Col span={8}>
              <Statistic
                title={
                  <>
                    Report Card GPA{' '}
                    <PopoutLink
                      url="https://go.canvascbl.com/help/gpas"
                      tracking={{
                        destinationName: destinationNames.helpdesk,
                        destinationType: destinationTypes.helpdesk.gpas,
                        via: vias.gpaReportCardQuestionIcon
                      }}
                    >
                      <Icon type="question-circle" />
                    </PopoutLink>
                  </>
                }
                value={roundNumberToDigits(gpa.unweighted.default, 2)}
              />
            </Col>
            <Col span={8}>
              <Statistic
                title={
                  <>
                    Traditional GPA{' '}
                    <PopoutLink
                      url="https://go.canvascbl.com/help/gpas"
                      tracking={{
                        destinationName: destinationNames.helpdesk,
                        destinationType: destinationTypes.helpdesk.gpas,
                        via: vias.gpaTraditionalQuestionIcon
                      }}
                    >
                      <Icon type="question-circle" />
                    </PopoutLink>
                  </>
                }
                value={roundNumberToDigits(gpa.unweighted.subgrades, 2)}
              />
            </Col>
          </Row>
          <Padding all={10} />
          <Typography.Text type="secondary">
            These unweighted GPAs only represent the current semester. Learn
            more{' '}
            <PopoutLink
              url="https://go.canvascbl.com/help/gpas"
              tracking={{
                destinationName: destinationNames.helpdesk,
                destinationType: destinationTypes.helpdesk.gpas,
                via: vias.gpaLearnMore
              }}
            >
              here
            </PopoutLink>
            .
          </Typography.Text>
        </>
      )}
      <Divider />
      <Typography.Text type="secondary">
        Please note that these grades may not be accurate or representative of
        your real grade. For the most accurate and up-to-date information,
        please consult someone from your school.
      </Typography.Text>
    </>
  );
}

const ConnectedGrades = connect((state) => ({
  courses: state.canvas.courses,
  plus: state.plus,
  gradedUsers: state.canvas.gradedUsers,
  grades: state.canvas.grades,
  users: state.canvas.users,
  activeUserId: state.canvas.activeUserId,
  user: state.canvas.user,
  observees: state.canvas.observees,
  showHiddenCourses: state.settings.showHiddenCourses,
  allGpas: state.canvas.gpa,
  error: state.error,
  loading: state.loading
}))(Grades);

export default ConnectedGrades;
