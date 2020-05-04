import moment from 'moment';

export default (courses, userId) =>
  courses.filter(
    (c) =>
      (c.end_at ? moment(c.end_at).isAfter(/* now */) : true) &&
      (userId
        ? c.enrollments.some(
            (e) => e.associated_user_id === userId || e.user_id === userId
          )
        : true)
  );
