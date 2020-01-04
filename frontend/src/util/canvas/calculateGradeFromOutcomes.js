const naGrade = { grade: 'N/A' };

export default (outcomeRollups, userId) => {
  const courseIds = Object.keys(outcomeRollups);

  const grades = {};

  courseIds.forEach(i => {
    const rollup = outcomeRollups[i];

    // array of outcomes for a specified class
    const scores = rollup[0].scores.map(s => s.score);

    if (scores.length < 1) {
      grades[i] = { grade: 'N/A' };
      return;
    }

    grades[i] = getGradeFromOutcomes(scores);
  });

  return grades;
};

function getGradeFromOutcomes(outcomes) {
  // outcomes, ex: [1,2,3,4]

  // what is 75% of outcomes.length
  const outcomesOverMinNeeded = Math.floor((75 * outcomes.length) / 100);

  // desc
  const sortedOutcomes = outcomes.sort((a, b) => b - a);
  // 75% of outcomes
  const countedOutcomes = sortedOutcomes.slice(0, outcomesOverMinNeeded);

  // lowest counted outcome
  const lowestCountedOutcome = countedOutcomes[countedOutcomes.length - 1];
  // lowest outcome (overall)
  const lowestOutcome = sortedOutcomes[sortedOutcomes.length - 1];

  // map lowest counted outcome to grade map
  for (let i = 0; i < gradeMap.length; i++) {
    // current available grade
    const grade = gradeMap[i][0];
    // not really max, but what all counted outcomes must be over
    const max = gradeMap[i][1][0];
    // what all outcomes must be over
    const min = gradeMap[i][1][1];

    // lowest outcome over minimum?
    if (lowestOutcome < min) {
      // one outcome is below min for this grade
      continue;
    }

    // counted outcomes must be above max, too
    if (lowestCountedOutcome < max) {
      // 75% of outcomes are not over max
      continue;
    }

    // student is eligible for this grade
    return {
      grade,
      sortedOutcomes,
      countedOutcomes,
      lowestCountedOutcome,
      lowestOutcome
    };
  }

  // we should never reach this point-- it would require a negative outcome
}

export const gradeMap = [
  ['A', [3.3, 3]],
  ['A-', [3.3, 2.5]],
  ['B+', [2.6, 2.5]],
  ['B', [2.6, 0]],
  ['B-', [2.6, 0]],
  ['C', [2.2, 0]],
  ['I', [0, 0]]
];

export const gradeMapByGrade = {
  A: { rank: 6, max: 3.3, min: 3 },
  'A-': { rank: 5, max: 3.3, min: 2.5 },
  'B+': { rank: 4, max: 2.6, min: 2.5 },
  B: { rank: 3, max: 2.6, min: 0 },
  'B-': { rank: 2, max: 2.6, min: 0 },
  C: { rank: 1, max: 2.2, min: 0 },
  I: { rank: 0, max: 0, min: 0 }
};
