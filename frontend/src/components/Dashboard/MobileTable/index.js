import React from 'react';
import * as PropTypes from 'prop-types';

import { Accordion as MobileAccordion, List as MobileList } from 'antd-mobile';

// Renders a table similar to antd's table, but built for antd-mobile.
function MobileTable(props) {
  const { columns, dataSource, accordionProps, indentChildren } = props;

  return (
    <MobileAccordion {...accordionProps}>
      {dataSource.map(d => (
        <MobileAccordion.Panel key={d.key} header={d.rowHeader}>
          <MobileList>
            {columns.map(c =>
              c.mobileRender ? (
                c.mobileRender(d[c.dataIndex], d)
              ) : c.render ? (
                c.render(d[c.dataIndex], d)
              ) : (
                <MobileList.Item
                  key={c.dataIndex}
                  extra={
                    c.extraRender
                      ? c.extraRender(d[c.dataIndex], d)
                      : d[c.dataIndex]
                  }
                  {...c.listItemProps}
                  style={{ paddingLeft: `${indentChildren ? 6 : 0}px` }}
                >
                  {c.title}
                </MobileList.Item>
              )
            )}
          </MobileList>
        </MobileAccordion.Panel>
      ))}
    </MobileAccordion>
  );
}

MobileTable.propTypes = {
  /** The columns of your table. Supports the title, dataIndex and render properties
   of the antd table. Also supports two custom props-- extraRender and listItemProps
   */
  columns: PropTypes.arrayOf(
    PropTypes.shape({
      // title of the column
      title: PropTypes.any.isRequired,
      // mobileTitle allows you to use the same columns object you use for desktop.
      // takes priority over title.
      mobileTitle: PropTypes.any,
      // how to find the data-- if we had { colA: 'fff' }, we would put
      // 'colA' here, because data['colA'] gives us the data.
      dataIndex: PropTypes.string.isRequired,
      // optional render function if you'd like to do something special with the row.
      // signature: (record) => <ReactComponent />
      // you may not supply render and extraRender-- render is expected to render both
      render: PropTypes.func,
      // mobileRender allows you to use the same columns object you use for desktop.
      // in your mobile table. takes priority over render.
      mobileRender: PropTypes.func,
      // extraRender lets you render just the extra on the List.Item.
      // same signature as antd table-- (text, record) => <ReactElement />
      // you may not supply render and extraRender-- render is expected to render both
      extraRender: PropTypes.func,
      // listItemProps lets you pass extra props to the List.Item rendered for every row.
      // you may not supply this with render-- render will replace the List.Item
      listItemProps: PropTypes.array
    })
  ).isRequired,
  // data contains the rows of your table
  dataSource: PropTypes.arrayOf(
    PropTypes.shape({
      // key is a unique identifier for every row and must be provided.
      key: PropTypes.string.isRequired,
      // rowHeader is the header of each table row
      rowHeader: PropTypes.any.isRequired
    })
  ).isRequired,
  // extra props to pass to the enclosing Accordion
  accordionProps: PropTypes.array,
  // whether to indent children by 5px
  indentChildren: PropTypes.bool
};

export default MobileTable;
