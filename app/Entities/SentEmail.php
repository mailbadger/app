<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class SentEmail extends Model implements Transformable
{
    use TransformableTrait;

    protected $table = 'sent_emails';

    protected $fillable = [
		'opens',
	];

    public function subscriber()
    {
        return $this->belongsTo('newsletters\Entities\Subscriber');
    }

    public function campaign()
    {
        return $this->belongsTo('newsletters\Entities\Campaign');
    }
}
