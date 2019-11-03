import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import v4 from 'uuid/v4';

import { Typography, Table, Icon, Spin, Tag, Skeleton, Popover } from 'antd';
import { Accordion as MobileAccordion, List as MobileList } from 'antd-mobile';

import {
  getUserCourses,
  getOutcomeRollupsForCourse
} from '../../../actions/canvas';

import calculateGradeFromOutcomes, {
  gradeMapByGrade
} from '../../../util/canvas/calculateGradeFromOutcomes';
import getActiveCourses from '../../../util/canvas/getActiveCourses';
import ErrorModal from '../ErrorModal';

import { ReactComponent as PopOutIcon } from '../../../assets/pop_out.svg';
import { ReactComponent as plusIcon } from '../../../assets/plus.svg';
import { desc } from '../../../util/stringSorter';
import PopoutLink from '../../PopoutLink';
import { getPreviousGrades } from '../../../actions/plus';
import moment from 'moment';
import { isMobile } from 'react-device-detect';

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
          url={`https://${localStorage.subdomain ||
            'canvas'}.instructure.com/courses/${record.id}`}
        >
          Open on Canvas <Icon component={PopOutIcon} />
        </PopoutLink>
      </div>
    )
  }
];

function Grades(props) {
  const [getCoursesId, setGetCoursesId] = useState('');
  const [
    getOutcomeRollupsForCourseIds,
    setGetOutcomeRollupsForCourseIds
  ] = useState([]);
  const [getPrevGradeId, setGetPrevGradeId] = useState('');

  const [loadingText, setLoadingText] = useState('');

  const allIds = [
    getCoursesId,
    getPrevGradeId,
    ...getOutcomeRollupsForCourseIds
  ];

  const {
    dispatch,
    token,
    subdomain,
    loading,
    error,
    user,
    courses,
    outcomeRollups,
    plus
  } = props;

  const err = error[Object.keys(error).filter(eid => allIds.includes(eid))[0]];

  const activeCourses = courses ? getActiveCourses(courses) : courses;

  useEffect(() => {
    if (allIds.some(id => loading.includes(id)) || err) {
      return;
    }
    if (!courses && !getCoursesId) {
      const id = v4();
      dispatch(getUserCourses(id, token, subdomain));
      setGetCoursesId(id);
      setLoadingText('your courses');
      return;
    }

    // if user AND no outcome rollups
    // or if we're missing a rollup for a class
    // and if we haven't fetched rollups already
    // fetch rollups
    if (
      (user && !outcomeRollups) ||
      ((() =>
        outcomeRollups
          ? activeCourses.some(c => !outcomeRollups[c.id])
          : false)() &&
        !getOutcomeRollupsForCourseIds.length)
    ) {
      const ids = [];
      activeCourses.forEach(c => {
        const id = v4();
        ids.push(id);
        dispatch(
          getOutcomeRollupsForCourse(id, user.id, c.id, token, subdomain)
        );
      });
      setGetOutcomeRollupsForCourseIds(ids);
      setLoadingText('your grades');
    }

    if (
      user &&
      plus.session &&
      plus.session.hasValidSubscription &&
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

  if (
    !user ||
    !plus.session ||
    !courses ||
    !outcomeRollups ||
    allIds.some(id => loading.includes(id))
  ) {
    return (
      <div align="center">
        <Spin />
        <span style={{ paddingTop: '20px' }} />
        <Typography.Title level={3}>
          {`Loading ${loadingText}...`}
        </Typography.Title>
      </div>
    );
  }

  const grades = calculateGradeFromOutcomes(outcomeRollups);

  const data = activeCourses.map(c => ({
    key: c.id,
    name: c.name,
    grade: grades[c.id] ? grades[c.id].grade : 'Error, try reloading',
    id: c.id,
    userHasValidSubscription: plus.session.hasValidSubscription,
    previousGrade: loading.includes(getPrevGradeId)
      ? 'loading'
      : plus &&
        plus.previousGrades &&
        !error[getPrevGradeId] &&
        plus.previousGrades.filter(pg => pg.courseId === c.id)[0] &&
        plus.previousGrades.filter(pg => pg.courseId === c.id)[0]
  }));

  if (isMobile) {
    return (
      <div>
        <Typography.Title level={2}>Grades</Typography.Title>
        <MobileAccordion>
          {data.map(d => (
            <MobileAccordion.Panel
              key={d.key}
              style={{ padding: '5px' }}
              header={
                <MobileList.Item extra={d.grade}>{d.name}</MobileList.Item>
              }
            >
              <MobileList>
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
      </div>
    );
  }

  return (
    <div>
      <Typography.Title level={2}>Grades</Typography.Title>
      <Typography.Text type="secondary">
        If you have a grade in a class, click on the name to see a detailed
        breakdown of your grade.
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
  outcomeRollups: state.canvas.outcomeRollups,
  user: state.canvas.user,
  token: state.canvas.token,
  subdomain: state.canvas.subdomain,
  error: state.error,
  loading: state.loading
}))(Grades);

export default ConnectedGrades;
