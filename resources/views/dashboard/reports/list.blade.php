@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/reports-list.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Reports</h1>
    <div id="reports"></div>
@endsection
