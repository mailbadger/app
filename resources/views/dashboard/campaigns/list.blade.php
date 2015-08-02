@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/campaigns-list.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">All campaigns</h1>
    <div id="campaigns"></div>
@endsection