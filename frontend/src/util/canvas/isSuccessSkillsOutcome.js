export default title => {
  const t = title.toLowerCase();

  return (
    t.startsWith('SS') ||
    (t.includes('teacher') && t.includes('evaluation')) ||
    (t.includes('teacher') && t.includes('assessment')) ||
    (t.includes('self') && t.includes('evaluation')) ||
    (t.includes('peer') && t.includes('evaluation')) ||
    (t.includes('peer') && t.includes('assessment')) ||
    (t.includes('self') && t.includes('evaluation'))
  );
};
