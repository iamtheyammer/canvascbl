import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import v4 from 'uuid/v4';

import { Typography, Table, Icon, Tag, Skeleton, Popover } from 'antd';
import { Accordion as MobileAccordion, List as MobileList } from 'antd-mobile';

import { gradeMapByGrade } from '../../../util/canvas/gradeMapByGrade';
import getActiveCourses from '../../../util/canvas/getActiveCourses';
import ErrorModal from '../ErrorModal';

import env from '../../../util/env';
import { ReactComponent as PopOutIcon } from '../../../assets/pop_out.svg';
import { ReactComponent as plusIcon } from '../../../assets/plus.svg';
import { desc } from '../../../util/sort';
import PopoutLink from '../../PopoutLink';
import { getPreviousGrades } from '../../../actions/plus';
import moment from 'moment';
import { isMobile } from 'react-device-detect';
import truncate from 'truncate';
import Loading from '../Loading';
import Padding from '../../Padding';
import ConnectedObserveeHandler from '../DashboardNav/ObserveeHandler';

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
      record.grade === 'N/A' || record.grade.toLowerCase().includes('error') ? (
        text
      ) : (
        <Link to={`/dashboard/grades/${record.id}`}>{text}</Link>
      )
  },
  {
    title: 'Grade',
    dataIndex: 'grade',
    key: 'grade',
    sorter: (a, b) => desc(a.grade, b.grade),
    defaultSortOrder: 'desc'
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
      <div>
        {record.grade !== 'N/A' &&
          !record.grade.toLowerCase().includes('error') && (
            <span>
              <Link to={`/dashboard/grades/${record.id}`}>See Breakdown</Link>
              {' | '}
            </span>
          )}
        <PopoutLink
          url={`https://${env.defaultSubdomain}.instructure.com/courses/${record.id}`}
        >
          Open on Canvas <Icon component={PopOutIcon} />
        </PopoutLink>
      </div>
    )
  }
];

function Grades(props) {
  const [getPrevGradeId, setGetPrevGradeId] = useState('');

  const {
    dispatch,
    grades,
    subdomain,
    loading,
    error,
    courses,
    users,
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

  if (err) {
    return <ErrorModal error={err} />;
  }

  if (!activeUser || !plus.session || !courses) {
    return <Loading text="grades" />;
  }

  const activeCourses =
    courses && activeUser ? getActiveCourses(courses, activeUser.id) : [];

  const previousGrades =
    plus &&
    plus.previousGrades &&
    plus.previousGrades.filter(pg => pg.canvasUserId === activeUserId);

  const data = activeCourses.map(c => ({
    key: c.id,
    name: c.name,
    grade: grades[activeUserId][c.id]
      ? grades[activeUserId][c.id].grade.grade
      : 'Error, try reloading',
    id: c.id,
    userHasValidSubscription: plus.session.has_valid_subscription,
    previousGrade: loading.includes(getPrevGradeId)
      ? 'loading'
      : plus &&
        plus.previousGrades &&
        !error[getPrevGradeId] &&
        previousGrades.filter(pg => pg.courseId === c.id)[0] &&
        previousGrades.filter(pg => pg.courseId === c.id)[0]
  }));

  const gradesTitle = (
    <Typography.Title level={2}>
      {observees && observees.length
        ? `${users[activeUserId].name.split(' ')[0]}'s Grades`
        : 'Grades'}
    </Typography.Title>
  );

  if (isMobile) {
    return (
      <div>
        {gradesTitle}
        <MobileAccordion>
          {data.map(d => (
            <MobileAccordion.Panel
              key={d.key}
              style={{ padding: '5px 5px 5px 0px' }}
              header={
                <div style={{ paddingRight: '6px' }}>
                  <div style={{ float: 'left', overflow: 'hidden' }}>
                    {truncate(d.name, 20)}
                  </div>
                  <div style={{ float: 'right' }}>{d.grade}</div>
                </div>
              }
            >
              <MobileList style={{ paddingLeft: '6px' }}>
                {d.grade !== 'N/A' && !d.grade.toLowerCase().includes('error') && (
                  <MobileList.Item>
                    <Link to={`/dashboard/grades/${d.id}`}>See Breakdown</Link>
                  </MobileList.Item>
                )}
                <MobileList.Item>
                  <PopoutLink
                    url={`https://${subdomain ||
                      'canvas'}.instructure.com/courses/${d.id}`}
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
              </MobileList>
            </MobileAccordion.Panel>
          ))}
        </MobileAccordion>
        {observees && observees.length > 0 && (
          <div>
            <Padding br />
            <Typography.Title level={3}>Switch Students</Typography.Title>
            <ConnectedObserveeHandler />
          </div>
        )}
      </div>
    );
  }

  return (
    <div>
      {gradesTitle}
      <Typography.Text type="secondary">
        If {observees && observees.length ? 'your student has' : 'you have'} a
        grade in a class, click on the name to see a detailed breakdown.
      </Typography.Text>
      <div style={{ marginBottom: '12px' }} />
      <Table columns={tableColumns} dataSource={data} />
      <Typography.Text type="secondary">
        Please note that these grades may not be accurate or representative of
        your real grade. For the most accurate and up-to-date information,
        please consult someone from your school.
      </Typography.Text>
    </div>
  );
}

const ConnectedGrades = connect(state => ({
  courses: state.canvas.courses,
  plus: state.plus,
  gradedUsers: state.canvas.gradedUsers,
  grades: state.canvas.grades,
  users: state.canvas.users,
  activeUserId: state.canvas.activeUserId,
  user: state.canvas.user,
  observees: state.canvas.observees,
  error: state.error,
  loading: state.loading
}))(Grades);

export default ConnectedGrades;
