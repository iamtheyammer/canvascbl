import React from 'react';
import * as PropTypes from 'prop-types';
import { Icon } from 'antd';
import { ReactComponent as PopOutIcon } from '../../assets/pop_out.svg';

function PopoutLink(props) {
  const { url, children, addIcon } = props;
  return (
    <a target="_blank" rel="noopener noreferrer" href={url}>
      {children}
      {addIcon && ' '}
      {addIcon && <Icon component={PopOutIcon} />}
    </a>
  );
}

PopoutLink.propTypes = {
  // url to link to
  url: PropTypes.string.isRequired,
  // whether to add a space and the PopOutIcon after your content (children)
  addIcon: PropTypes.bool
};

export default PopoutLink;
