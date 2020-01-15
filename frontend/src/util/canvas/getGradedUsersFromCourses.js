import { flatten } from 'lodash';

export default courses => [
  // using [...someSet] coerces the Set into an Array
  // using a Set because it doesn't allow duplicates
  ...new Set(
    // we will get an array of graded users from each class, we want a single array
    flatten(
      // for each course
      courses.map(c =>
        // get all enrollments (associated_user_id is available for observers, otherwise it's user_id)
        c.enrollments.map(e => e.associated_user_id || e.user_id)
      )
    )
  )
];
