@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/components/templates/templates.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">All templates</h1>
    <div id="templates"></div>
@endsection