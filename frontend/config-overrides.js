const {
  override,
  fixBabelImports,
  addWebpackPlugin
} = require('customize-cra');
const Dotenv = require('dotenv-webpack');

module.exports = override(
  fixBabelImports('import-antd', {
    libraryName: 'antd',
    libraryDirectory: 'es',
    style: 'css'
  }),
  fixBabelImports('import-mobile', {
    libraryName: 'antd-mobile',
    style: 'css'
  }),
  fixBabelImports('import-lodash', {
    libraryName: 'lodash',
    libraryDirectory: '',
    camel2DashComponentName: false // default: true
  }),
  addWebpackPlugin(new Dotenv())
);
