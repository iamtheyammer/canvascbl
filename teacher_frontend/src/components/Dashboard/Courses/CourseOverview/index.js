import React, { Fragment, useEffect } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import styled from "styled-components";
import {
  Button,
  Divider,
  Icon,
  Input,
  Select,
  Skeleton,
  Table,
  Typography
} from "antd";
import Padding from "../../../Padding";
import {
  getCourseEnrollments,
  getDistanceLearningGradesOverview
} from "../../../../actions/canvas";
import moment from "moment";
import {
  clearFilters,
  filterName,
  filterNameType
} from "../../../../actions/filters";
import { desc } from "../../../../util/sort";
import PopoutLink from "../../../PopoutLink";

const gradeOverviewTableColumns = [
  {
    title: "Student Name",
    dataIndex: "studentName",
    key: "studentName",
    sorter: (a, b) => desc(a.studentName, b.studentName)
  },
  {
    title: "Grade",
    dataIndex: "grade",
    key: "grade",
    render: text => <Typography.Text strong>{text}</Typography.Text>
  },
  {
    title: "Timestamp",
    dataIndex: "timestamp",
    key: "timestamp",
    render: text => moment.utc(text).calendar()
  },
  {
    title: "Actions",
    dataIndex: "actions",
    key: "actions",
    render: items =>
      items.map((a, i) => (
        <Fragment key={i}>
          {a}
          {i !== items.length - 1 && <Divider type="vertical" />}
        </Fragment>
      ))
  }
];

const HalfWidthSkeleton = styled(Skeleton)`
  width: 50%;
`;

function CourseOverview(props) {
  const {
    distanceLearningPairs,
    loadingDlGradesOverviews,
    getDlGradesOverviewError,
    dlGradesOverviews,
    loadingEnrollments,
    getCourseEnrollmentsError,
    enrollments,
    filters,
    match,
    dispatch
  } = props;

  const courseId = match.params.courseId;
  const dlPair =
    distanceLearningPairs &&
    distanceLearningPairs.filter(
      p =>
        courseId === `${p.original_course_id}_${p.distance_learning_course_id}`
    )[0];

  const overview = dlGradesOverviews && dlGradesOverviews[courseId];
  const dlEnrolls =
    enrollments && dlPair && enrollments[dlPair.distance_learning_course_id];

  useEffect(() => {
    if (
      dlPair &&
      !loadingDlGradesOverviews &&
      !getDlGradesOverviewError &&
      !overview
    ) {
      dispatch(
        getDistanceLearningGradesOverview(
          dlPair.original_course_id,
          dlPair.distance_learning_course_id
        )
      );
    }
  }, [
    dlPair,
    loadingDlGradesOverviews,
    getDlGradesOverviewError,
    overview,
    dispatch
  ]);

  useEffect(() => {
    if (
      dlPair &&
      !loadingEnrollments &&
      !getCourseEnrollmentsError &&
      !dlEnrolls
    ) {
      dispatch(getCourseEnrollments(dlPair.distance_learning_course_id));
    }
  }, [
    dlPair,
    loadingEnrollments,
    getCourseEnrollmentsError,
    dlEnrolls,
    dispatch
  ]);

  if (!dlPair && distanceLearningPairs) {
    return <Redirect to="/dashboard/courses" />;
  }

  if (getDlGradesOverviewError) {
    return (
      <Typography.Text type="danger">
        We encountered an error getting grades. Please try again later or
        contact support.
      </Typography.Text>
    );
  }

  if (getCourseEnrollmentsError) {
    return (
      <Typography.Text type="danger">
        We encountered an error getting students. Please try again later or
        contact support.
      </Typography.Text>
    );
  }

  // reflect filters, add filter changer
  const tableData =
    overview &&
    overview
      .filter(o => {
        const enroll =
          dlEnrolls && dlEnrolls.filter(e => e.user_id === o.user_id)[0];

        const name =
          enroll && enroll.user.name && enroll.user.name.toLowerCase();
        const fname = filters.name && filters.name.toLowerCase();

        if (filters.name) {
          switch (filters.nameType) {
            case "includes":
              if (!name.includes(fname)) {
                return false;
              }
              break;
            case "startsWith":
              if (!name.startsWith(fname)) {
                return false;
              }
              break;
            case "endsWith":
              if (!name.endsWith(fname)) {
                return false;
              }
              break;
            default:
              if (!name.includes(fname)) {
                return false;
              }
              break;
          }
        }

        return true;
      })
      .map((o, i) => {
        const enroll =
          dlEnrolls && dlEnrolls.filter(e => e.user_id === o.user_id)[0];

        return {
          studentName: enroll ? (
            enroll.user.name
          ) : (
            <HalfWidthSkeleton paragraph={false} />
          ),
          grade: o.grade.grade,
          timestamp: o.timestamp,
          user: enroll && enroll.user,
          actions: [
            enroll && (
              <PopoutLink url={`mailto:${enroll.user.login_id}`}>
                <Icon type="mail" /> Email Student
              </PopoutLink>
            )
          ],
          key: i
        };
      });

  return (
    <>
      <Typography.Title level={2}>
        {dlPair ? (
          <>Grade Overview for {dlPair.course_name}</>
        ) : (
          <Skeleton paragraph={false} />
        )}
      </Typography.Title>
      <Typography.Text type="secondary">
        See grades for every student in your class.
        <br />
        <br />
        Missing grades? It can take up to three hours for grades to sync to
        CanvasCBL.
      </Typography.Text>
      <Typography.Title level={3}>Filters</Typography.Title>
      <Typography.Text>
        Use filters to control which grades you see.
      </Typography.Text>
      <Padding all={5} />
      <Button
        onClick={() => dispatch(clearFilters())}
        icon="delete"
        type="danger"
        size="small"
      >
        Clear Filters
      </Button>
      <Padding all={5} />
      <Typography.Title level={4}>Name</Typography.Title>
      <Input
        disabled={!dlEnrolls}
        placeholder="Joe Smith"
        value={filters.name}
        onChange={e => dispatch(filterName(e.target.value))}
        style={{ width: "40%" }}
        addonBefore={
          <Select
            value={filters.nameType || "includes"}
            onChange={type => dispatch(filterNameType(type))}
            style={{ width: 125 }}
            disabled={!dlEnrolls}
          >
            <Select.Option value="includes">Includes</Select.Option>
            <Select.Option value="startsWith">Starts With</Select.Option>
            <Select.Option value="endsWith">Ends With</Select.Option>
          </Select>
        }
      />
      <Padding all={10} />
      <Table
        columns={gradeOverviewTableColumns}
        dataSource={tableData}
        loading={!tableData}
      />
    </>
  );
}

const ConnectedCourseOverview = connect(state => ({
  distanceLearningPairs: state.canvas.distanceLearningPairs,
  loadingDlGradesOverviews: state.canvas.loadingDistanceLearningGradesOverview,
  getDlGradesOverviewError: state.canvas.getDistanceLearningGradesOverviewError,
  dlGradesOverviews: state.canvas.distanceLearningGradesOverviews,
  loadingEnrollments: state.canvas.courseEnrollmentsAreLoading,
  getCourseEnrollmentsError: state.canvas.getCourseEnrollmentsError,
  enrollments: state.canvas.enrollments,
  filters: state.filters
}))(CourseOverview);

export default ConnectedCourseOverview;
