import React from 'react';
import { Icon, Result, Typography } from 'antd';

function ComingSoon(props) {
  return (
    <>
      <Result
        icon={<Icon type="bell" theme="twoTone" />}
        title="What do you want to see here?"
        extra={
          <>
            <Typography.Text>
              What do you want from CanvasCBL for teachers? Let us know via the
              Provide Feedback button in the top right.
            </Typography.Text>
          </>
        }
      />
    </>
  );
}

export default ComingSoon;
