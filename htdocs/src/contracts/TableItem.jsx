import React from 'react';
import { formatAmount } from '../helpers.js';

class TableItem extends React.Component {
  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e) {
    e.preventDefault();
    this.props.onClick( this.props.contract );
  }

  render() {
    let release = this.props.contract.releases[0];
    return (
      <tr>
        <td width="50%">
          <p className="lbl">{release.tender.id}</p>
          <a onClick={this.handleClick}>{release.tender.description}</a>
        </td>
        <td>
          <p className="lbl">Monto</p>
          <h4>{formatAmount( release.tender.value.amount || 0 )}</h4>
          <p className="lbl">Proveedor</p>
          <p>{release.awards[0].suppliers[0].name}</p>
        </td>
        <td>
          <p className="lbl">Comprador</p>
          <p>{release.buyer.name}</p>
        </td>
        <td width="15%">
          <a className="btn-black active" onClick={this.handleClick}>Ver Contrato</a>
        </td>
      </tr>
    );
  }
}

export default TableItem;
