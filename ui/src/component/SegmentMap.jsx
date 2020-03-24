import React, {Component} from 'react'
import {Map, Marker, TileLayer, Tooltip} from "react-leaflet";
import axios from "axios";
import {Dimmer, Loader, Segment} from "semantic-ui-react";

export default class SegmentMap extends Component {

    constructor(props) {
        super(props);
        this.mapRef = React.createRef();
        this.state = {
            lat: 48.50113756,
            lng: 9.744529724,
            zoom: 14,
            rankings: [],
            loader: false,
        }
    }

    onClickDone = (evt) => {
        let bounds = this.mapRef.current.leafletElement.getBounds();
        this.setState({
            viewport: this.props.viewport,
            loader: true,
        });

        const comp = this;
        axios.post('/api/v1/rankings', {
            sw: {lat: bounds.getSouthWest().lat, lng: bounds.getSouthWest().lng},
            ne: {lat: bounds.getNorthEast().lat, lng: bounds.getNorthEast().lng},
        })
            .then(function (response) {
                comp.setState({
                    rankings: response.data,
                    loader: false,
                })
            })
            .catch(function (error) {
                console.error(error);
            });
    };

    render() {
        const position = [this.state.lat, this.state.lng];
        let i = 0;
        const markers = this.state.rankings.map(s => <Marker key={i++} position={[s.start.lat, s.start.lng]}>
            <Tooltip permanent>
                {s.segmentName} <br/> {s.position}/{s.all}
            </Tooltip>
        </Marker>);
        return (
            <Segment>
                <Dimmer active={this.state.loader} inverted>
                    <Loader size='large'>Loading</Loader>
                </Dimmer>
                <Map ref={this.mapRef} onClick={this.onClickDone} center={position} zoom={this.state.zoom}>
                    <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                    />
                    {markers}
                </Map>
            </Segment>


        )
    }
}
