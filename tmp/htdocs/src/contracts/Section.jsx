import React from 'react';
import { getParameter, formatAmount, formatDate } from '../helpers.js';
import { Description } from '../general.jsx';
import TableItem from './TableItem.jsx';
import SearchResults from './SearchResults.jsx';
import SearchBar from './SearchBar.jsx';
import Details from './Details.jsx';

class Section extends React.Component {
  constructor(props) {
    super(props);
    this.runQuery = this.runQuery.bind(this);
    this.showDetails = this.showDetails.bind(this);
    this.hideDetails = this.hideDetails.bind(this);
    this.state = {
      items: [],
      details: null
    };
  }

  runQuery(query) {
    let url = '/query/gacm';
    if( getParameter('bucket') ) {
      url = '/query/' + getParameter('bucket');
    }

    // Submit query and update component state with results
    $.ajax({
      type: "POST",
      url: url,
      data: {
        query: JSON.stringify(query)
      },
      success: (function( res ) {
        this.setState({
          items: JSON.parse( res )
        });
      }).bind(this)
    });
  }

  showDetails(contract) {
    this.setState({
      details: contract
    });
  }

  hideDetails() {
    this.setState({
      details: null
    });
  }

  render() {
    if( !this.state.details ) {
      return (
        <div>
          <Description
            title="Contratos"
            color="blue"
            content="Consulta cada contrato que está registrado en Testigo Social 2.0. Podrás encontrar información para cada una de las etapas del procedimiento de contratación, desde su planeación hasta su implementación." />
          <SearchBar onSubmit={this.runQuery} />
          <SearchResults items={this.state.items} onSelection={this.showDetails} />
        </div>
      );
    } else {
      return (
        <div>
          <Description
            title="Contratos"
            color="blue"
            content="En esta sección podrás visualizar la información completa y agregada de un contrato, desde su planeación hasta su implementación. Como base se utiliza el Estándar de Datos para las Contrataciones Abiertas." />
          <Details contract={this.state.details} onClose={this.hideDetails} />
        </div>
      );
    }
  }
}

export default Section;
