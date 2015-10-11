@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/settings.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Settings</h1>
    <div class="row" id="settings"></div>
@endsection