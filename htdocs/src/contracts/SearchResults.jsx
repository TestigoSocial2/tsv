import React from 'react';
import TableItem from './TableItem.jsx';

class SearchResults extends React.Component {
  constructor(props) {
    super(props);
    this.itemSelected = this.itemSelected.bind(this);
  }

  itemSelected(contract) {
    this.props.onSelection(contract);
  }

  render() {
    let content = (<tr><td><h4><i>Sin resultados que mostrar...</i></h4></td></tr>);
    if( this.props.items.length > 0 ) {
      content = this.props.items.map((contract) =>
        <TableItem onClick={this.itemSelected} key={contract.releases[0].id} contract={contract} />
      );
    }

    return(
      <div className="inner-row">
        <div className="row">
          <div className="col-md-12">
            <h2 className="mobile-result">Resultados de la busqueda</h2>
            <p className="mobile-result">Estos son los procedimientos de contratación registrados en Testigo Social 2.0 que satisfacen los creiterios de selección proporcionados</p>
            <table className="table table-striped latest-contracts">
              <tbody>
                {content}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }
}

export default SearchResults;
