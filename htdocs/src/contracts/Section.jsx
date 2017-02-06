import React from 'react';
import Description from '../base/Description.jsx';
import TableItem from './TableItem.jsx';
import SearchResults from './SearchResults.jsx';
import SearchBar from './SearchBar.jsx';
import Details from './Details.jsx';
import { getParameter } from '../helpers.js';

class Section extends React.Component {
  constructor(props) {
    super(props);
    this.runQuery = this.runQuery.bind(this);
    this.state = {
      items: [],
      details: null
    };
  }

  // Submit query and update component state with results
  runQuery(query) {
    let url = '/query/gacm';
    if( getParameter('bucket') ) {
      url = '/query/' + getParameter('bucket');
    }

    $.ajax({
      type: "POST",
      url: url,
      data: {
        query: JSON.stringify(query)
      },
      success: (res) => this.setState({ items: JSON.parse( res ) })
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
          <SearchResults items={this.state.items} onSelection={(c) => this.setState({details:c})} />
        </div>
      );
    } else {
      return (
        <div>
          <Description
            title="Contratos"
            color="blue"
            content="En esta sección podrás visualizar la información completa y agregada de un contrato, desde su planeación hasta su implementación. Como base se utiliza el Estándar de Datos para las Contrataciones Abiertas." />
          <Details contract={this.state.details} onClose={() => this.setState({details:null})} />
        </div>
      );
    }
  }
}

export default Section;
