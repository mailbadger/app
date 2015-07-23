@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/components/campaigns-table.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">All campaigns</h1>
    <div class="row">
        <div class="col-lg-4">
            <a href="{{url('dashboard/new-campaign')}}" class="btn btn-success btn-lg">
                <span class="glyphicon glyphicon-plus"></span> Create new campaign
            </a>
        </div>
    </div>

    <div class="row">
        <div class="col-lg-12" id="campaigns"></div>
    </div>
@endsection