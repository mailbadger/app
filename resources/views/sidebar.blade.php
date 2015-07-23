<div class="col-sm-3 col-md-2 sidebar">
    <ul class="nav nav-sidebar">
        <li class="@if($activeSidebar == 'dashboard') active @endif">
            <a href="{{url('dashboard')}}">All Campaigns</a>
        </li>
        <li class="@if($activeSidebar == 'new-campaign') active @endif">
            <a href="{{url('dashboard/new-campaign')}}">Create new campaign</a>
        </li>
    </ul>
    <ul class="nav nav-sidebar">
        <li class="@if($activeSidebar == 'templates') active @endif">
            <a href="{{url('dashboard/templates')}}">Templates</a>
        </li>
        <li class="@if($activeSidebar == 'new-template') active @endif">
            <a href="{{url('dashboard/new-template')}}">Create new template</a>
        </li>
    </ul>
    <ul class="nav nav-sidebar">
        <li class="@if($activeSidebar == 'sub-lists') active @endif">
            <a href="{{url('dashboard/subscribers')}}">Subscriber lists</a>
        </li>
        <li class="@if($activeSidebar == 'new-subs') active @endif">
            <a href="{{url('dashboard/new-subscribers')}}">Create new list</a>
        </li>
    </ul>
    <ul class="nav nav-sidebar">
        <li class="@if($activeSidebar == 'reports') active @endif">
            <a href="{{url('dashboard/reports')}}">Reports</a>
        </li>
    </ul>
</div>