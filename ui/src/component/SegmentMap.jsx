import React, {Component} from 'react'
import {Map, Marker, TileLayer, Tooltip} from "react-leaflet";

export default class SegmentMap extends Component {

    constructor(props) {
        super(props);
        this.mapRef = React.createRef();
        this.state = {
            lat: 48.50113756,
            lng: 9.744529724,
            zoom: 14,
        }
    }

    onClickDone = (evt) => {
        let bounds = this.mapRef.current.leafletElement.getBounds();
        console.log(bounds);
        this.setState({
            viewport: this.props.viewport,
        })
    };

    render() {
        const position = [this.state.lat, this.state.lng];
        const segments = [
            {lat: 48.471508, lng: 10.081507},
            {lat: 48.670443, lng: 10.079522}];

        let i = 0;
        const markers = segments.map(s => <Marker key={i++} position={s}>
            <Tooltip permanent>
                A pretty CSS3 popup. <br/> Easily customizable.
            </Tooltip>
        </Marker>);
        return (
            <Map ref={this.mapRef} onClick={this.onClickDone} center={position} zoom={this.state.zoom}>
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                />
                {markers}
            </Map>
        )
    }
}
