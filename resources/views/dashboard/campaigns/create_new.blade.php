@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/components/campaign-form.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Create new campaign</h1>
    <div class="row" id="new-campaign"></div>
@endsection